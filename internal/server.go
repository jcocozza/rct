package internal

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
)

type Server struct {
	// Where the server will listen
	Addr string
	// the token that clients will have to send to validate
	// if a client does not send the correct token, then no data will be received and the connection will be closed
	token         string
	tokenRequired bool
	clipboard     Clipboard
	results       chan string
}

func NewServer(addr string, tkn string, results chan string) *Server {
	tknRequired := tkn != ""
	return &Server{
		Addr:          addr,
		token:         tkn,
		tokenRequired: tknRequired,
		clipboard:     &clipper{},
		results:       results,
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
	//t := []byte(strings.ReplaceAll(string(msgBuf), `\n`, "\n"))
	//err = s.clipboard.Write(t)
	err = s.clipboard.Write(msgBuf)
	if err != nil {
		_ = respondError(conn, fmt.Errorf("failed to write to clipboard: %w", err))
		return
	}
	if s.results != nil {
		s.results <- string(msgBuf)
	}
}

func (s *Server) isAlive() bool {
	conn, err := net.Dial("tcp", s.Addr)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func (s *Server) Run() error {
	if s.isAlive() {
		return fmt.Errorf("server is already running")
	}
	listen, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to start listening: %s", err.Error())
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			// TODO: handle these nicely
			continue
		}
		s.handleConnection(conn)
	}
}

// call the executable in its own process
// return pid if success. -1 if failed to start
func (s *Server) RunDetached() (int, error) {
	if s.isAlive() {
		return -1, fmt.Errorf("server is already running")
	}
	exe, err := os.Executable()
	if err != nil {
		return -1, err
	}
	cmd := exec.Command(exe, "listen")

	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return -1, err
	}
	return cmd.Process.Pid, nil
}
