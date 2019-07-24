package utils

import (
	"encoding/binary"
	"math"
)

type ByteOrder struct {
	binary.ByteOrder
}

var LittleEndian = ByteOrder{ByteOrder: binary.LittleEndian}
var BigEndian = ByteOrder{ByteOrder: binary.BigEndian}

func (o *ByteOrder) UInt16ToBytes(v uint16) []byte {
	buf := [2]byte{}
	o.PutUint16(buf[:], v)
	return buf[:]
}

func (o *ByteOrder) UInt32ToBytes(v uint32) []byte {
	buf := [4]byte{}
	o.PutUint32(buf[:], v)
	return buf[:]
}

func (o *ByteOrder) UInt64ToBytes(v uint64) []byte {
	buf := [4]byte{}
	o.PutUint64(buf[:], v)
	return buf[:]
}

func (o *ByteOrder) Float32ToBytes(v float32) []byte {
	buf := [4]byte{}
	o.PutUint32(buf[:], math.Float32bits(v))
	return buf[:]
}
