package engin

import (
	"net/netip"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHConfig holds the configuration for SSH connection
type SSHConfig struct {
	User     string
	Password string
	Host     string
	Port     uint16

	addr netip.AddrPort
}

// connect establishes an SSH connection
func (s SSHConfig) connect() (*ssh.Client, error) {
	if s.addr.Addr().IsLoopback() {
		return nil, nil
	}
	sshConfig := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", s.addr.String(), sshConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}
