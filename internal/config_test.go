package internal

import (
	"testing"
)

func TestHostValidate(t *testing.T) {
	cases := []struct {
		name       string
		host       Host
		expectsErr bool
	}{
		{
			name:       "Valid host with IP and port",
			host:       Host{Addr: "127.0.0.1:8080"},
			expectsErr: false,
		},
		{
			name:       "Valid host with hostname and port",
			host:       Host{Addr: "example.com:9090"},
			expectsErr: false,
		},
		{
			name:       "Invalid host with missing port",
			host:       Host{Addr: "127.0.0.1"},
			expectsErr: true,
		},
		{
			name:       "Invalid host with non-numeric port",
			host:       Host{Addr: "127.0.0.1:port"},
			expectsErr: true,
		},
		{
			name:       "Invalid host with out-of-range port",
			host:       Host{Addr: "127.0.0.1:70000"},
			expectsErr: true,
		},
		{
			name:       "Invalid empty address",
			host:       Host{Addr: ""},
			expectsErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.host.Validate()
			if (err != nil) && !tc.expectsErr {
				t.Errorf("expected error: %v, got: %v", tc.expectsErr, err)
			}
		})
	}
}

func TestRCTConfigValidate(t *testing.T) {
	cases := []struct {
		name       string
		config     RCTConfig
		expectsErr bool
	}{
		{
			name: "Valid config with server and delivery",
			config: RCTConfig{
				Server:   Host{Addr: "127.0.0.1:8080"},
				Delivery: []Host{{Addr: "example.com:9090"}},
			},
			expectsErr: false,
		},
		{
			name: "Valid config with server only",
			config: RCTConfig{
				Server: Host{Addr: "127.0.0.1:8080"},
			},
			expectsErr: false,
		},
		{
			name: "Valid config with delivery only",
			config: RCTConfig{
				Delivery: []Host{{Addr: "example.com:9090"}},
			},
			expectsErr: false,
		},
		{
			name:       "Invalid config with empty server and delivery",
			config:     RCTConfig{},
			expectsErr: true,
		},
		{
			name: "Invalid config with invalid server",
			config: RCTConfig{
				Server: Host{Addr: "127.0.0.1"},
			},
			expectsErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if (err != nil) && !tc.expectsErr {
				t.Errorf("expected error: %v, got: %v", tc.expectsErr, err)
			}
		})
	}
}
