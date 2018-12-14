package utils

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"time"
)

//SSHConfig type to map SSH Configs
type SSHConfig struct {
	User       string
	Pass       string
	Host       string
	Ciphers    []string
	AuthMethod string
	PrivateKey string
}

// SSHClient dials up our target device.
func SSHClient(c SSHConfig) (*ssh.Client, error) {
	// Set up SSH auth method
	var sshAuth ssh.AuthMethod
	switch c.AuthMethod {
	case "KeyboardInteractive":
		sshAuth = ssh.KeyboardInteractive(c.sshInteractive)
		break
	case "PublicKeys": // TODO: Test this
		sshAuth = c.sshPublicKeys()
		break
	// default to password auth
	default:
		sshAuth = ssh.Password(c.Pass)
		break
	}

	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            []ssh.AuthMethod{sshAuth},
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
		return nil, err
	}

	return client, err
}

// sshInteractive sets up KeyboardInteractive ssh
func (c SSHConfig) sshInteractive(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	answers = make([]string, len(questions))
	for n := range questions {
		answers[n] = c.Pass
	}

	return answers, nil
}

// sshPublicKeys sets up Public Key auth for ssh
func (c SSHConfig) sshPublicKeys() ssh.AuthMethod {
	key, err := ioutil.ReadFile(c.PrivateKey)
	privateKey, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Println("Unable to parse SSH private key, reverting to password authentication")
		return ssh.Password(c.Pass)
	}
	return ssh.PublicKeys(privateKey)
}
