package srp

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"log"
	"math/big"
	"strings"
	"xcore/utils"
)

var (
	errInvalidPassHash = errors.New("invalid password hash")
)

const (
	saltSizeBits = 32 * 8
)

type SRP struct {
	g, xN      *big.Int
	b, xB      *big.Int
	xA, xK, xM *big.Int

	salt, verifier *big.Int
}

func newSRP() *SRP {
	return &SRP{
		xN: NewBigIntWithHex("894B645E89E1535BBDAD5B8B290650530801B18EBFBF5E8FAB3C82872A3E9BB7"),
		g:  NewBigIntWithHex("7"),
	}
}

func NewSRPWithPasswordHash(ph string) (*SRP, error) {
	srp := newSRP()

	passHashBytes, err := hex.DecodeString(ph)
	if err != nil || len(passHashBytes) != sha1.Size {
		return nil, errInvalidPassHash
	}

	srp.initWithPasswordHash(passHashBytes)
	return srp, nil
}

func NewSRPWithPasswordKey(token string) *SRP {
	srp := newSRP()
	t := strings.Split(token, ":")

	if len(t) != 2 {
		panic("invalid token")
	}

	srp.verifier = NewBigIntWithHex(t[0])
	srp.salt = NewBigIntWithHex(t[1])
	return srp
}

func (srp *SRP) initWithPasswordHash(passHash []byte) {
	srp.salt = RandBigInt(saltSizeBits)
	h := sha1.New()
	hashMustWrite(h, utils.ReversedBytes(srp.salt.Bytes()))

	ph := passHash
	phLen := len(ph)
	if phLen < sha1.Size {
		left := make([]uint8, sha1.Size-phLen)
		right := make([]uint8, sha1.Size)
		ph = append(left, right...)
	}
	hashMustWrite(h, passHash)

	x := newBigIntFromBytes(utils.ReversedBytes(h.Sum(nil)))
	// verifier = (g ^ x) % N
	srp.verifier = (&(big.Int{})).Exp(srp.g, x, srp.xN)
}

func (srp *SRP) GetVerifier() *big.Int {
	return srp.verifier
}

func (srp *SRP) GetSalt() *big.Int {
	return srp.salt
}

//noinspection GoSnakeCaseUsage
func (srp *SRP) GetGenerator() *big.Int {
	return srp.g
}

func (srp *SRP) GetPrime() *big.Int {
	return srp.xN
}

func (srp *SRP) GetEphemeralKey() *big.Int {
	if srp.xB == nil {
		srp.b = RandBigInt(19 * 8)
		// B=(k*v + g^b % N) % N
		vMul := (&(big.Int{})).Mul(srp.verifier, big.NewInt(3))
		gMod := (&(big.Int{})).Exp(srp.g, srp.b, srp.xN)
		sum := (&(big.Int{})).Add(vMul, gMod)
		srp.xB = (&(big.Int{})).Mod(sum, srp.xN)

		if srp.xB.BitLen() < 32 * 8 {
			data := make([]uint8, 32)
			copy(data, srp.xB.Bytes())
			srp.xB.SetBytes(data)
		}
	}
	return srp.xB
}

func (srp *SRP) GetProof() []byte {
	p := sha1.New()
	p.Write(utils.ReversedBytes(srp.xA.Bytes()))
	p.Write(utils.ReversedBytes(srp.xM.Bytes()))
	p.Write(utils.ReversedBytes(srp.xK.Bytes()))
	return p.Sum(nil)
}

func (srp *SRP) GetPublicKey() *big.Int {
	return srp.xK
}

func (srp *SRP) ValidateClientProof(accName string, clientProof []byte, clientPubKey []byte) bool {
	srp.xA = newBigIntFromBytes(utils.ReversedBytes(clientPubKey))
	zero := big.NewInt(0)
	AModN := (&(big.Int{})).Mod(srp.xA, srp.xN)
	if AModN.Cmp(zero) == 0 {
		return false
	}

	sha := sha1.New()
	hashMustWrite(sha, utils.ReversedBytes(srp.xA.Bytes()))
	hashMustWrite(sha, utils.ReversedBytes(srp.xB.Bytes()))

	// (A * (v.ModExp(u, N))).ModExp(b, N);
	uBytes := sha.Sum(nil)
	u := newBigIntFromBytes(utils.ReversedBytes(uBytes))

	tmp0 := (&(big.Int{})).Exp(srp.verifier, u, srp.xN)
	tmp1 := (&(big.Int{})).Mul(srp.xA, tmp0)
	S := (&(big.Int{})).Exp(tmp1, srp.b, srp.xN)

	if S.BitLen() < 16 * 8 {
		data := make([]uint8, 16)
		copy(data, S.Bytes())
		S.SetBytes(data)
	}

	tmpArr0 := utils.ReversedBytes(S.Bytes())
	tmpArr1 := [16]uint8{}
	tmpK := [40]uint8{}

	for i := 0; i < 16; i++ {
		tmpArr1[i] = tmpArr0[i*2]
	}

	sha.Reset()
	hashMustWrite(sha, tmpArr1[:])

	t1Hash := sha.Sum(nil)
	for i := 0; i < 20; i++ {
		tmpK[i*2] = t1Hash[i]
	}

	for i := 0; i < 16; i++ {
		tmpArr1[i] = tmpArr0[i*2+1]
	}

	sha.Reset()
	hashMustWrite(sha, tmpArr1[:])
	t1Hash = sha.Sum(nil)

	for i := 0; i < 20; i++ {
		tmpK[i*2+1] = t1Hash[i]
	}
	srp.xK = newBigIntFromBytes(utils.ReversedBytes(tmpK[:]))

	sha.Reset()
	hashMustWrite(sha, utils.ReversedBytes(srp.xN.Bytes()))
	hsh := sha.Sum(nil)

	sha.Reset()
	hashMustWrite(sha, utils.ReversedBytes(srp.g.Bytes()))
	gHash := sha.Sum(nil)

	for i := 0; i < sha1.Size; i++ {
		hsh[i] ^= gHash[i]
	}

	sha.Reset()
	hashMustWrite(sha, []byte(accName))
	accNameHash := sha.Sum(nil)

	sha.Reset()
	hashMustWrite(sha, hsh)
	hashMustWrite(sha, accNameHash)
	hashMustWrite(sha, utils.ReversedBytes(srp.salt.Bytes()))
	hashMustWrite(sha, utils.ReversedBytes(srp.xA.Bytes()))
	hashMustWrite(sha, utils.ReversedBytes(srp.xB.Bytes()))
	hashMustWrite(sha, utils.ReversedBytes(srp.xK.Bytes()))

	expectedProof := sha.Sum(nil)
	if subtle.ConstantTimeCompare(expectedProof, clientProof) == 1 {
		srp.xM = newBigIntFromBytes(utils.ReversedBytes(expectedProof))
		return true
	} else {
		return false
	}
}

func hashMustWrite(h hash.Hash, p []byte) {
	if _, err := h.Write(p); err != nil {
		log.Panic(err)
	}
}

func NewBigIntWithHex(h string) *big.Int {
	i, ok := (&(big.Int{})).SetString(h, 16)
	if !ok {
		panic("Invalid hex")
	}
	return i
}

func newBigIntFromBytes(b []byte) *big.Int {
	return (&(big.Int{})).SetBytes(b)
}
func RandBigInt(bits int) *big.Int {
	n := bits / 8
	if bits%8 != 0 {
		n += 1
	}
	return newBigIntFromBytes(randomBytes(n))
}

func randomBytes(count int) []byte {
	b := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic("Random source is broken!")
	}
	return b
}

//func copyBigInt(i *big.Int) *big.Int {
//	return (&(big.Int{})).Set(i)
//}
