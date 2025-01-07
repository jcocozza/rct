package internal

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

func resolve(hostname string) (string, error) {
	ips, err := net.LookupHost(hostname)
	if err != nil {
		return "", err
	}
	return ips[0], nil // TODO: improve this logic
}

func handleRemote(addr string) (string, error) {
	var finalRemoteAddr string
	isValidIP := net.ParseIP(addr) != nil
	if isValidIP {
		finalRemoteAddr = addr
	} else {
		resolvedIP, err := resolve(addr)
		if err != nil {
			return "", err
		}
		finalRemoteAddr = resolvedIP
	}
	return finalRemoteAddr, nil
}

func scpSendFile(ctx context.Context, username string, password string, remoteIP string, fileName string, file io.Reader) error {
	ip, err := handleRemote(remoteIP)
	if err != nil {
		return err
	}
	// Create an SCP client
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConn, err := ssh.Dial("tcp4", fmt.Sprintf("%s:22", ip), sshConfig)
	if err != nil {
		return err
	}
	client, err := scp.NewClientBySSH(sshConn)
	if err != nil {
		return err
	}
	defer client.Close()

	// Copy the file to the remote server
	err = client.CopyFile(ctx, file, fileName, "0644")
	if err != nil {
		return fmt.Errorf("Failed to copy file: %v", err)
	}
	return nil
}
