package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const CFG_PATH = ".rct.json"

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

type RCTConfig struct {
	// the host that will be started to listen for incoming
	Server Host `json:"server"`
	// tcp servers to try and send data to
	Delivery []Host `json:"delivery"`
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
