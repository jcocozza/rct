package internal

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type Server struct {
	// Where the server will listen
	Addr  string
	// the token that clients will have to send to validate
	// if a client does not send the correct token, then no data will be received and the connection will be closed
	token string
	tokenRequired bool
	clipboard Clipboard
}

func NewServer(addr string, tkn string) *Server {
	tknRequired := tkn != ""
	return &Server{
		Addr: addr,
		token: tkn,
		tokenRequired: tknRequired,
		clipboard: &clipper{},
	}
}

func respondError(conn net.Conn, err error) error {
	_, e := conn.Write([]byte(err.Error()))
	return e
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	if s.tokenRequired {
		// read token
		var tknLength int32
		err := binary.Read(reader, binary.BigEndian, &tknLength)
		if err != nil {
			_ = respondError(conn, fmt.Errorf("unable to read token length"))
			return
		}
		tknBuf := make([]byte, tknLength)
		_, err = io.ReadFull(reader, tknBuf)
		if err != nil {
			_ = respondError(conn, fmt.Errorf("unable to read token"))
			return
		}
		token := string(tknBuf)

		// token validation
		if token != s.token {
			_ = respondError(conn, fmt.Errorf("invalid token"))
			return
		}
	}

	// read message
	var msgLength int32
	err := binary.Read(reader, binary.BigEndian, &msgLength)
	if err != nil {
		_ = respondError(conn, fmt.Errorf("unable to read message length"))
		return
	}
	msgBuf := make([]byte, msgLength)
	_, err = io.ReadFull(reader, msgBuf)
	if err != nil {
		_ = respondError(conn, fmt.Errorf("unable to read message"))
		return
	}
	// process message
	err = s.clipboard.Write(msgBuf)
	if err != nil {
		_ = respondError(conn, fmt.Errorf("failed to write to clipboard: %w", err))
		return
	}
}

func (s *Server) Run() {
	listen, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		s.handleConnection(conn)
	}
}
