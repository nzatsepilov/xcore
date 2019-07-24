package net

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
	"xcore/utils"
)

type Socket struct {
	conn     *net.TCPConn
	readBuf  *utils.Buffer
	writeBuf *utils.Buffer

	isClosed bool
	onClose  func(err error)
}

func NewSocket(conn *net.TCPConn) *Socket {
	return &Socket{
		conn:     conn,
		readBuf:  utils.NewBuffer(),
		writeBuf: utils.NewBuffer(),
	}
}

func (s *Socket) RemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func (s *Socket) ReceiveData() error {
	deadline := time.Now().Add(time.Second * 5)
	if err := s.conn.SetReadDeadline(deadline); err != nil {
		return err
	}

	count, err := s.readBuf.ReadFrom(s.conn)
	if err != nil {
		return err
	}

	if count == 0 {
		return io.EOF
	}

	return nil
}

func (s *Socket) ReceiveDataRecursive(expectedSize int) error {
	for s.ReadBufferSize() < expectedSize {
		if err := s.ReceiveData(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Socket) WriteByte(v byte) error {
	return s.writeBuf.WriteByte(v)
}

func (s *Socket) WriteUInt16(v uint16) error {
	return s.WriteBytes(utils.LittleEndian.UInt16ToBytes(v))
}

func (s *Socket) WriteUInt32(v uint32) error {
	return s.WriteBytes(utils.LittleEndian.UInt32ToBytes(v))
}

func (s *Socket) WriteUInt64(v uint64) error {
	return s.WriteBytes(utils.LittleEndian.UInt64ToBytes(v))
}

func (s *Socket) WriteFloat32(v float32) error {
	return s.WriteBytes(utils.LittleEndian.Float32ToBytes(v))
}

func (s *Socket) WriteBool(v bool) error {
	if v {
		return s.WriteByte(1)
	} else {
		return s.WriteByte(0)
	}
}

func (s *Socket) MustWriteByte(v byte) *Socket {
	if err := s.writeBuf.WriteByte(v); err != nil {
		log.Panic(err)
	}
	return s
}

func (s *Socket) MustWriteUInt16(v uint16) *Socket {
	if err := s.WriteUInt16(v); err != nil {
		log.Panic(err)
	}
	return s
}

func (s *Socket) MustWriteUInt32(v uint32) *Socket {
	if err := s.WriteUInt32(v); err != nil {
		log.Panic(err)
	}
	return s
}

func (s *Socket) MustWriteFloat32(v float32) *Socket {
	if err := s.WriteFloat32(v); err != nil {
		log.Panic(err)
	}
	return s
}

func (s *Socket) MustWriteBool(v bool) *Socket {
	if err := s.WriteBool(v); err != nil {
		log.Panic(err)
	}
	return s
}

func (s *Socket) WriteBytes(buf []byte) error {
	_, err := s.writeBuf.Write(buf)
	return err
}

func (s *Socket) MustWriteBytes(buf []byte) *Socket {
	if err := s.WriteBytes(buf); err != nil {
		log.Panic(err)
	}
	return s
}

func (s *Socket) BeginWrite() *Socket {
	s.writeBuf.Reset()
	return s
}

func (s *Socket) CommitWrite() error {
	deadline := time.Now().Add(time.Second * 5)
	if err := s.conn.SetWriteDeadline(deadline); err != nil {
		return err
	}

	_, err := s.writeBuf.WriteTo(s.conn)
	if s.writeBuf.Len() > 0 {
		return io.ErrShortWrite
	}

	return err
}

func (s *Socket) MustCommitWrite() {
	if err := s.CommitWrite(); err != nil {
		log.Panic(err)
	}
}

func (s *Socket) ReadByte() (byte, error) {
	return s.readBuf.ReadByte()
}

func (s *Socket) ReadByteTo(v *byte) error {
	var err error
	*v, err = s.ReadByte()
	return err
}

func (s *Socket) MustReadByte() byte {
	b, err := s.readBuf.ReadByte()
	if err != nil {
		log.Panic(err)
	}
	return b
}

func (s *Socket) MustReadByteTo(v *byte) *Socket {
	if err := s.ReadByteTo(v); err != nil {
		log.Panic(err)
	}
	return s
}

func (s *Socket) ReadBytes(count int) ([]byte, error) {
	return s.readBuf.ReadBytes(count)
}

func (s *Socket) ReadBytesTo(v *[]byte, count int) error {
	var err error
	*v, err = s.ReadBytes(count)
	return err
}

func (s *Socket) MustReadBytes(count int) []byte {
	b, err := s.ReadBytes(count)
	if err != nil {
		log.Panic(err)
	}
	return b
}

func (s *Socket) ReadBytesWithDelimiter(delim byte) ([]byte, error) {
	return s.readBuf.ReadBytesWithDelimiter(delim)
}

func (s *Socket) MustReadBytesWithDelimiter(delim byte) []byte {
	b, err := s.ReadBytesWithDelimiter(delim)
	if err != nil {
		log.Panic(err)
	}
	return b
}

func (s *Socket) ReadString() (string, error) {
	b, err := s.ReadBytesWithDelimiter(0)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *Socket) MustReadString() string {
	str, err := s.ReadString()
	if err != nil {
		log.Panic(err)
	}
	return str
}

func (s *Socket) MustReadBytesTo(v *[]byte, count int) *Socket {
	if err := s.ReadBytesTo(v, count); err != nil {
		log.Panic(err)
	}
	return s
}

func (s *Socket) ReadTo(v interface{}) error {
	return binary.Read(s.readBuf, binary.LittleEndian, v)
}

func (s *Socket) MustReadTo(v interface{}) {
	if err := s.ReadTo(v); err != nil {
		log.Panic(err)
	}
}

func (s *Socket) ReadUInt16() (uint16, error) {
	var v uint16
	err := s.ReadTo(v)
	return v, err
}

func (s *Socket) MustReadUInt16() uint16 {
	var v uint16
	s.MustReadTo(v)
	return v
}

func (s *Socket) ReadUInt32() (uint32, error) {
	var v uint32
	err := s.ReadTo(v)
	return v, err
}

func (s *Socket) MustReadUInt32() uint32 {
	var v uint32
	s.MustReadTo(v)
	return v
}

func (s *Socket) ReadUInt64() (uint64, error) {
	var v uint64
	err := s.ReadTo(v)
	return v, err
}

func (s *Socket) MustReadUInt64() uint64 {
	var v uint64
	s.MustReadTo(v)
	return v
}

func (s *Socket) ReadBool() (bool, error) {
	var v bool
	err := s.ReadTo(v)
	return v, err
}

func (s *Socket) MustReadBool() bool {
	var v bool
	s.MustReadTo(&v)
	return v
}

func (s *Socket) ReadFloat32() (float32, error) {
	var v float32
	err := s.ReadTo(v)
	return v, err
}

func (s *Socket) MustReadFloat32() float32 {
	var v float32
	s.MustReadTo(v)
	return v
}

func (s *Socket) ReadFloat64() (float64, error) {
	var v float64
	err := s.ReadTo(v)
	return v, err
}

func (s *Socket) MustReadFloat64() float64 {
	var v float64
	s.MustReadTo(v)
	return v
}

func (s *Socket) ReadBufferSize() int {
	return s.readBuf.Len()
}

func (s *Socket) WriteBufferSize() int {
	return s.writeBuf.Len()
}

func (s *Socket) Close() error {
	if s.isClosed {
		return nil
	}

	err := s.conn.Close()
	s.isClosed = true

	if s.onClose != nil {
		s.onClose(err)
	}
	return err
}

func (s *Socket) OnClose(handler func(err error)) *Socket {
	s.onClose = handler
	return s
}
