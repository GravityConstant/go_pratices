package self_rpc

import (
	"encoding/binary"
	"io"
	"net"
)

type Session struct {
	conn net.Conn
}

func NewSession(conn net.Conn) *Session {
	return &Session{conn: conn}
}

// write
func (s *Session) Write(data []byte) error {
	// head + body
	buf := make([]byte, 4+len(data))
	// write head
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	// copy body
	copy(buf[4:], data)
	_, err := s.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// read
func (s *Session) Read() ([]byte, error) {
	// head
	header := make([]byte, 4)
	// read body length
	_, err := io.ReadFull(s.conn, header)
	if err != nil {
		return nil, err
	}
	// read data
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(s.conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
