package auth

import (
	"crypto/sha1"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	goNet "net"
	"strconv"
	"strings"
	"xcore/core/models"
	"xcore/core/net"
	"xcore/core/srp"
	"xcore/utils"
)

var (
	errUnexpectedOpcode = errors.New("unexpected auth opcode")
)

type sessionStatus uint8

const (
	logonChallengeStatus sessionStatus = iota
	logonProofStatus
	authorizedStatus
	reconnectProofStatus
	closedStatus
)

const (
	realmListMsgSize      = 4
	reconnectProofMsgSize = 57
)

type sessionHandler struct {
	status  sessionStatus
	msgSize int
	handler func(s *session) error
}

var sessionHandlers map[opcode]*sessionHandler
var versionChallenge []uint8

type Session interface {
	authorize()
}

type session struct {
	sock *net.Socket

	id      string
	status  sessionStatus
	account *models.Account

	accRepo    AccountRepository
	realmList  *realmProvider
	srp        *srp.SRP
	reconProof *big.Int
}

func init() {
	sessionHandlers = map[opcode]*sessionHandler{
		logonChallengeOpcode: {
			status:  logonChallengeStatus,
			msgSize: logonChallengeSize,
			handler: (*session).handleLogonChallengeOpcode,
		},
		logonProofOpcode: {
			status:  logonProofStatus,
			msgSize: logonProofSize,
			handler: (*session).handleLogonProofOpcode,
		},
		realmlistOpcode: {
			status:  authorizedStatus,
			msgSize: realmListMsgSize,
			handler: (*session).handleRealmListOpcode,
		},
		reconnectChallengeOpcode: {
			status:  logonChallengeStatus,
			msgSize: logonChallengeSize,
			handler: (*session).handleReconnectChallengeOpcode,
		},
		reconnectProofOpcode: {
			status:  reconnectProofStatus,
			msgSize: reconnectProofMsgSize,
			handler: (*session).handleReconnectProofOpcode,
		},
	}
	versionChallenge = []uint8{0xBA, 0xA3, 0x1E, 0x99, 0xA0, 0x0B, 0x21, 0x57, 0xFC, 0x37, 0x3F, 0xB3, 0x69, 0xCD, 0xD2, 0xF1}
}

func newSession(id string, conn *goNet.TCPConn, accRepo AccountRepository, rs *realmProvider) Session {
	sock := net.NewSocket(conn)
	sock.OnClose(func(err error) {
		if err != nil {
			log.Printf("Auth session [%v] closed with error: %v", id, err)
		} else {
			log.Printf("Auth session [%v] closed", id)
		}
	})
	return &session{
		sock:      sock,
		id:        id,
		accRepo:   accRepo,
		realmList: rs,
	}
}

func (s *session) initSRP() error {
	k := s.account.PasswordKey.String
	if s.account.PasswordKey.Valid && len(k) > 0 {
		s.srp = srp.NewSRPWithPasswordKey(k)
		return nil
	}

	var err error = nil
	s.srp, err = srp.NewSRPWithPasswordHash(s.account.PasswordHash)
	if err != nil {
		return err
	}

	if err = s.createAndSaveAuthToken(); err != nil {
		return err
	}

	return nil
}

func (s *session) authorize() {
	log.Printf("Auth session [%v] started (%v)", s.id, s.sock.RemoteAddr())

	if err := s.continueAuth(); err != nil {
		if err := s.sock.Close(); err != nil {
			log.Panic(err)
		}
	}
}

func (s *session) continueAuth() error {
	err := s.sock.ReceiveData()
	if err == io.EOF {
		s.status = closedStatus
		return s.sock.Close()
	}

	opRaw, err := s.sock.ReadByte()
	if err == io.EOF {
		s.status = closedStatus
		return nil
	}

	op := opcode(opRaw)
	h := sessionHandlers[op]
	if h == nil || h.status != s.status {
		s.status = closedStatus
		log.Printf("Received unexpected opcode: %v", op)
		return errUnexpectedOpcode
	}

	msgSize := h.msgSize
	recvSize := s.sock.ReadBufferSize()

	if recvSize < msgSize {
		log.Printf("Received malformed packed with %d size, but expected %d size", recvSize+1, msgSize)
		s.status = closedStatus
		return s.sock.Close()
	}

	log.Printf("Handling opcode %v\n", op)
	if err := h.handler(s); err != nil {
		s.status = closedStatus
		return err
	}

	return nil
}

func (s *session) handleLogonChallengeOpcode() error {
	buf, err := s.sock.ReadBytes(logonChallengeSize)
	if err != nil {
		return err
	}

	challenge, err := newLogonChallenge(buf)
	if err != nil {
		return err
	}

	accLen := int(challenge.accNameLen)
	if s.sock.ReadBufferSize() < accLen {
		return s.sock.Close()
	}

	accBuf, err := s.sock.ReadBytes(accLen)
	if err != nil {
		return err
	}

	return s.handleLogonChallenge(challenge, string(accBuf))
}

func (s *session) handleLogonProofOpcode() error {
	buf, err := s.sock.ReadBytes(logonProofSize)
	if err != nil {
		return err
	}

	p, err := newLogonProof(buf)
	if err != nil {
		return err
	}

	return s.handleLogonProof(p)
}

func (s *session) handleRealmListOpcode() error {
	_, err := s.sock.ReadBytes(realmListMsgSize)
	if err != nil {
		return err
	}

	count := uint16(s.realmList.GetRealmsCount())

	realmsBuf := utils.NewBuffer()
	realmsBuf.MustWriteBytes(utils.LittleEndian.UInt32ToBytes(0))
	realmsBuf.MustWriteBytes(utils.LittleEndian.UInt16ToBytes(count))

	for i := 0; i < int(count); i++ {
		r := s.realmList.GetRealm(i)
		realmsBuf.MustWriteByte(byte(r.Type))

		if r.IsLocked {
			realmsBuf.MustWriteByte(1)
		} else {
			realmsBuf.MustWriteByte(0)
		}

		realmsBuf.MustWriteByte(byte(r.Flag))

		realmsBuf.MustWriteBytes([]byte(r.Name))
		realmsBuf.MustWriteByte(0)

		realmsBuf.MustWriteBytes([]byte(r.Address))
		realmsBuf.MustWriteByte(0)

		population := float32(r.Population)
		realmsBuf.MustWriteBytes(utils.LittleEndian.Float32ToBytes(population))

		realmsBuf.MustWriteByte(r.CharactersCount)
		realmsBuf.MustWriteByte(byte(r.Timezone))

		realmsBuf.MustWriteByte(r.ID)

		if r.Flag.Has(models.RealmFlagSpecifyBuild) {
			v := strings.Split(r.Version, ".")
			var err error
			vMajor, err := strconv.Atoi(v[0])
			vMinor, err := strconv.Atoi(v[1])
			vBugFix, err := strconv.Atoi(v[2])
			vBuild, err := strconv.Atoi(v[3])

			if err != nil {
				return err
			}

			realmsBuf.MustWriteByte(byte(vMajor))
			realmsBuf.MustWriteByte(byte(vMinor))
			realmsBuf.MustWriteByte(byte(vBugFix))
			realmsBuf.MustWriteBytes(utils.LittleEndian.UInt16ToBytes(uint16(vBuild)))
		}
	}

	// Unused
	realmsBuf.MustWriteBytes(utils.LittleEndian.UInt16ToBytes(0x0010))

	s.sock.BeginWrite()
	s.sock.MustWriteByte(byte(realmlistOpcode))
	s.sock.MustWriteUInt16(uint16(realmsBuf.Len()))
	s.sock.MustWriteBytes(realmsBuf.Bytes())

	if err := s.sock.CommitWrite(); err != nil {
		return err
	}

	return s.continueAuth()
}

func (s *session) handleReconnectChallengeOpcode() error {
	buf, err := s.sock.ReadBytes(logonChallengeSize)
	if err != nil {
		return err
	}

	challenge, err := newLogonChallenge(buf)
	if err != nil {
		return err
	}

	accLen := int(challenge.accNameLen)
	if s.sock.ReadBufferSize() < accLen {
		return s.sock.Close()
	}

	accBuf, err := s.sock.ReadBytes(accLen)
	if err != nil {
		return err
	}

	return s.handleReconnectChallenge(challenge, string(accBuf))
}

func (s *session) handleLogonChallenge(payload *logonChallenge, accountName string) error {
	acc, err := s.accRepo.GetAccountWithName(accountName)
	if err != nil {
		return err
	}

	if acc == nil {
		return s.closeWithResult(resultUnknownAccount, logonChallengeOpcode)
	}

	s.account = acc

	if err := s.initSRP(); err != nil {
		return err
	}

	B := utils.ReversedBytes(s.srp.GetEphemeralKey().Bytes())
	g := utils.ReversedBytes(s.srp.GetGenerator().Bytes())
	N := utils.ReversedBytes(s.srp.GetPrime().Bytes())
	salt := utils.ReversedBytes(s.srp.GetSalt().Bytes())

	if len(B) != 32 {
		log.Panicln("Invalid public key size")
	}

	if len(g) != 1 {
		log.Panicln("Invalid generator size")
	}

	if len(N) != 32 {
		log.Panicln("Invalid prime size")
	}

	if len(salt) != 32 {
		log.Panicln("Invalid salt size")
	}

	s.sock.BeginWrite()

	s.sock.MustWriteByte(byte(logonChallengeOpcode))
	s.sock.MustWriteByte(0x00)
	s.sock.MustWriteByte(byte(resultSuccess))
	s.sock.MustWriteBytes(B)
	s.sock.MustWriteByte(0x01)
	s.sock.MustWriteBytes(g)
	s.sock.MustWriteByte(32)
	s.sock.MustWriteBytes(N)
	s.sock.MustWriteBytes(salt)
	s.sock.MustWriteBytes(versionChallenge)
	s.sock.MustWriteByte(0) // security flags

	if err := s.sock.CommitWrite(); err != nil {
		return err
	}

	s.status = logonProofStatus
	return s.continueAuth()
}

func (s *session) handleLogonProof(logonProof *logonProof) error {
	accName := strings.ToUpper(s.account.Name)
	if !s.srp.ValidateClientProof(accName, logonProof.xM1[:], logonProof.xA[:]) {
		s.sock.BeginWrite()
		s.sock.MustWriteByte(byte(logonProofOpcode))
		s.sock.MustWriteByte(byte(resultUnknownAccount))
		s.sock.MustWriteByte(3)
		s.sock.MustWriteByte(0)

		if err := s.sock.CommitWrite(); err != nil {
			return err
		}

		return s.continueAuth()
	}

	s.account.SessionKey = sql.NullString{
		String: s.srp.GetPublicKey().Text(16),
		Valid:  true,
	}
	if err := s.accRepo.SaveAccount(s.account); err != nil {
		return err
	}

	s.sock.BeginWrite()
	s.sock.MustWriteByte(byte(logonProofOpcode))
	s.sock.MustWriteByte(0) // error
	s.sock.MustWriteBytes(s.srp.GetProof())

	s.sock.MustWriteUInt32(uint32(accountFlagPropass))
	s.sock.MustWriteUInt32(0) // survey id
	s.sock.MustWriteUInt16(0) // login flags

	if err := s.sock.CommitWrite(); err != nil {
		return err
	}

	s.status = authorizedStatus
	return s.continueAuth()
}

func (s *session) handleReconnectChallenge(challenge *logonChallenge, accName string) error {
	acc, err := s.accRepo.GetAccountWithName(accName)
	if err != nil {
		return err
	}

	if acc == nil {
		return s.closeWithResult(resultUnknownAccount, reconnectChallengeOpcode)
	}

	if !acc.SessionKey.Valid {
		return s.closeWithResult(resultSessionExpired, reconnectChallengeOpcode)
	}

	s.account = acc
	s.reconProof = srp.RandBigInt(16 * 8)

	s.sock.BeginWrite()
	s.sock.MustWriteByte(byte(reconnectChallengeOpcode))
	s.sock.MustWriteByte(0)
	s.sock.MustWriteBytes(utils.ReversedBytes(s.reconProof.Bytes()))

	// 16 bytes of zeros
	unk := [16]uint8{}
	s.sock.MustWriteBytes(unk[:])

	if err := s.sock.CommitWrite(); err != nil {
		return err
	}

	s.status = reconnectProofStatus
	return s.continueAuth()
}

func (s *session) handleReconnectProofOpcode() error {
	R1 := s.sock.MustReadBytes(16)
	R2 := s.sock.MustReadBytes(20)
	_ = s.sock.MustReadBytes(20) // R3 (unused)
	_ = s.sock.MustReadByte()    // number of keys (unused)

	failure := func() error {
		s.sock.BeginWrite()
		s.sock.MustWriteByte(byte(reconnectProofOpcode))
		s.sock.MustWriteByte(byte(resultUnknownAccount))
		s.sock.MustWriteByte(3)
		s.sock.MustWriteByte(0)

		if err := s.sock.CommitWrite(); err != nil {
			return err
		}

		return s.continueAuth()
	}

	if !s.account.SessionKey.Valid {
		return failure()
	}

	K := srp.NewBigIntWithHex(s.account.SessionKey.String)

	h := sha1.New()
	h.Write([]byte(strings.ToUpper(s.account.Name)))
	h.Write(R1)
	h.Write(utils.ReversedBytes(s.reconProof.Bytes()))
	h.Write(utils.ReversedBytes(K.Bytes()))
	expectedR2 := h.Sum(nil)

	if subtle.ConstantTimeCompare(expectedR2, R2) == 0 {
		return failure()
	}

	s.sock.BeginWrite().
		MustWriteByte(byte(reconnectProofOpcode)).
		MustWriteByte(0).
		MustWriteUInt16(0)

	if err := s.sock.CommitWrite(); err != nil {
		return err
	}

	s.status = authorizedStatus

	return s.continueAuth()
}

func (s *session) closeWithResult(result result, command opcode) error {
	s.status = closedStatus

	s.sock.BeginWrite().
		MustWriteByte(byte(command)).
		MustWriteByte(0x00).
		MustWriteByte(byte(result))

	if err := s.sock.CommitWrite(); err != nil {
		return err
	}

	return s.continueAuth()
}

func (s *session) createAndSaveAuthToken() error {
	verifier := s.srp.GetVerifier().Text(16)
	salt := s.srp.GetSalt().Text(16)

	s.account.PasswordKey = sql.NullString{
		Valid:  true,
		String: fmt.Sprintf("%v:%v", verifier, salt),
	}

	if err := s.accRepo.SaveAccount(s.account); err != nil {
		return err
	}
	return nil
}
