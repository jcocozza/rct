package internal

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const CFG_PATH = ".rct.json"
const DEFAULT_PORT = "54321"
const timeout = 3 * time.Second


/*
{
	"server": {
		"addr": "",
		"token": ""
	},
	"delivery": [
		{
			"addr": "",
			"token": ""
		}
	]
}
*/

type Host struct {
	// address of the tcp server
	Addr string `json:"addr"`
	// token that will validate data on the tcp server
	Token string `json:"token"`
}

func (h Host) Validate() error {
	host, port, err := net.SplitHostPort(h.Addr)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("invalid port: %s", port)
	}
	if net.ParseIP(host) == nil && len(host) == 0 {
		return fmt.Errorf("invalid server: %s", host)
	}
	return nil
}

type RCTConfig struct {
	// the host that will be started to listen for incoming
	Server Host `json:"server"`
	// tcp servers to try and send data to
	Delivery []Host `json:"delivery"`
}


func (c RCTConfig) Validate() error {
	if strings.TrimSpace(c.Server.Addr) != "" {
		err := c.Server.Validate()
		if err != nil {
			return err
		}
	}
	if len(c.Delivery) > 0 {
		for _, hst := range c.Delivery {
			err := hst.Validate()
			if err != nil {
				return err
			}
		}
	}
	if strings.TrimSpace(c.Server.Addr) == "" && len(c.Delivery) == 0 {
		return fmt.Errorf("config must contain either server, delivery or both")
	}
	return nil
}

func readConfig(path string) (RCTConfig, error) {
	fbytes, err := os.ReadFile(path)
	if err != nil {
		return RCTConfig{}, err
	}
	var sc RCTConfig
	err = json.Unmarshal(fbytes, &sc)
	if err != nil {
		return RCTConfig{}, err
	}
	err = sc.Validate()
	if err != nil {
		return RCTConfig{}, err
	}
	return sc, nil
}

func ReadConfig() (RCTConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return RCTConfig{}, err
	}
	configPath := filepath.Join(homeDir, CFG_PATH)
	return readConfig(configPath)
}


func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddress := conn.LocalAddr().(*net.UDPAddr)
	return localAddress.IP.String(), nil
}

func GenerateLocal(port string) (RCTConfig, error) {
	ip, err := getLocalIP()
	if err != nil {
		return RCTConfig{}, err
	}
	addr := net.JoinHostPort(ip, port)
	return RCTConfig{
		Server: Host{Addr: addr},
		Delivery: []Host{},
	}, nil
}

func GenerateRemote(port string) (RCTConfig, error) {
	ip, err := getLocalIP()
	if err != nil {
		return RCTConfig{}, err
	}
	addr := net.JoinHostPort(ip, port)
	return RCTConfig{
		Server: Host{},
		Delivery: []Host{
			{Addr: addr},
		},
	}, nil
}
