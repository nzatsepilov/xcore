package auth

import (
	"bytes"
	"encoding/binary"
)

/*
typedef struct AUTH_LOGON_PROOF_C
{
    uint8_t   cmd;
    uint8_t   A[32];
    uint8_t   M1[20];
    uint8_t   crc_hash[20];
    uint8_t   number_of_keys;
    uint8_t   securityFlags;
} sAuthLogonProof_C;
*/

const logonProofSize = 74

type logonProof struct {
	xA            [32]uint8
	xM1           [20]uint8
	crcHash       [20]uint8
	keysCount     uint8
	securityFlags uint8
}

func newLogonProof(b []byte) (*logonProof, error) {
	buf := bytes.NewBuffer(b)
	p := new(logonProof)

	if err := binary.Read(buf, binary.LittleEndian, &p.xA); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p.xM1); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p.crcHash); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p.keysCount); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &p.securityFlags); err != nil {
		return nil, err
	}

	return p, nil
}

type serverLogonProofPayload struct {
}
