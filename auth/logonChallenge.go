package auth

import (
	"bytes"
	"encoding/binary"
	"xcore/utils"
)

const logonChallengeSize = 33

type logonChallenge struct {
	error        uint8
	size         uint16
	gameName     string
	version      [3]uint8
	build        uint16
	platform     string
	os           string
	country      string
	timezoneBias uint32 // minutes from UTC
	ip           [4]uint8
	accNameLen   uint8
}

func newLogonChallenge(b []byte) (*logonChallenge, error) {
	buf := bytes.NewBuffer(b)
	c := new(logonChallenge)

	readString := func(str *string, size int, reverse bool) error {
		a := make([]uint8, size)
		if err := binary.Read(buf, binary.LittleEndian, &a); err != nil {
			return err
		}
		if reverse {
			utils.ReverseBytes(a)
		}
		*str = string(a)
		return nil
	}

	if err := binary.Read(buf, binary.LittleEndian, &c.error); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &c.size); err != nil {
		return nil, err
	}
	if err := readString(&c.gameName, 4, true); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &c.version); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &c.build); err != nil {
		return nil, err
	}
	if err := readString(&c.platform, 4, true); err != nil {
		return nil, err
	}
	if err := readString(&c.os, 4, true); err != nil {
		return nil, err
	}
	if err := readString(&c.country, 4, true); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &c.timezoneBias); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &c.ip); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &c.accNameLen); err != nil {
		return nil, err
	}
	return c, nil
}
