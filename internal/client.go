package internal

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

type Client struct {
	// the address of the server that client will send data to
	ServerAddr string
	// the token that the server will validate before it processes the data
	serverToken string
}

func NewClient(serverAddr string, serverToken string) *Client {
	return &Client{
		ServerAddr:  serverAddr,
		serverToken: serverToken,
	}
}

// Send the txt to the server
//
// will send the following (in order):
// 1. the token length
// 2. the token
// 3. the message length
// 3. the message
func (c *Client) Send(txt string) error {
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if c.serverToken != "" {
		// send api key
		tokenLen := int32(len(c.serverToken))
		err = binary.Write(conn, binary.BigEndian, tokenLen)
		if err != nil {
			return err
		}
		_, err = conn.Write([]byte(c.serverToken))
		if err != nil {
			return err
		}
	}

	// send message
	msgLen := int32(len(txt))
	err = binary.Write(conn, binary.BigEndian, msgLen)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(txt))
	if err != nil {
		return err
	}

	errBuf := make([]byte, 1024)
	n, err := conn.Read(errBuf)
	if err != nil && err != io.EOF {
		return err
	}
	if n != 0 {
		return fmt.Errorf("%s", string(errBuf))
	}
	return nil
}

// sendToHost sends the text to a single host and prints errors if verbose mode is enabled.
func sendToHost(host Host, txt string) error {
	client := NewClient(host.Addr, host.Token)
	return client.Send(txt)
}

// deliver sends the text to all hosts concurrently and reports any errors.
func deliver(hosts []Host, txt string) []error {
	var wg sync.WaitGroup
	errLst := make([]error, len(hosts))
	for i, host := range hosts {
		wg.Add(1)
		go func(host Host) {
			defer wg.Done()
			if err := sendToHost(host, txt); err != nil {
				errLst[i] = err
			}
		}(host)
	}
	wg.Wait()
	return errLst
}
