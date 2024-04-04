package config

import (
	"errors"
	"net"
	"os"
)

const (
	authClientHost = "AUTH_HOST"
	authClientPort = "AUTH_PORT"
)

type authConfig struct {
	host string
	port string
}

// NewAuthConfig - создает конфиг gRPC
func NewAuthConfig() (AuthConfig, error) {
	host := os.Getenv(authClientHost)
	if len(host) == 0 {
		return nil, errors.New("auth host not found in environments")
	}

	port := os.Getenv(authClientPort)
	if len(host) == 0 {
		return nil, errors.New("auth port not found in environments")
	}

	return authConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg authConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
