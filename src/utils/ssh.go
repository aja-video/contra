package utils

import (
	"golang.org/x/crypto/ssh"
	"time"
)

//SSHConfig type to map SSH Configs
type SSHConfig struct {
	User     string
	Password string
	Host     string
	Ciphers  []string
}

// SSHClient dials up our target device.
func SSHClient(c SSHConfig) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: this should be an option
		Timeout:         time.Second * 10,            //TODO: this should be an option
	}
	if c.Ciphers != nil {
		config.Config = ssh.Config{
			Ciphers: c.Ciphers,
		}

	}

	client, err := ssh.Dial("tcp", c.Host, config)

	if err != nil {
		panic(err)
	}

	return client, err
}
