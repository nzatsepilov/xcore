package auth

import (
	"bytes"
	"encoding/binary"
)

const reconnectProofSize = 57

type reconnectProof struct {
	xR1       [16]uint8
	xR2       [20]uint8
	xR3       [20]uint8
	keysCount uint8
}

func newReconnectProof(b []byte) (*reconnectProof, error) {
	buf := bytes.NewBuffer(b)
	c := new(reconnectProof)

	if err := binary.Read(buf, binary.LittleEndian, &c.xR1); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &c.xR2); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &c.xR3); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &c.keysCount); err != nil {
		return nil, err
	}

	return c, nil
}
