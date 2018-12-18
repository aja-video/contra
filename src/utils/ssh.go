package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"time"
)

//SSHConfig type to map SSH Configs
type SSHConfig struct {
	User          string
	Pass          string
	Host          string
	Ciphers       []string
	AuthMethod    string
	PrivateKey    string
	AllowInsecure bool
	SSHTimeout    time.Duration
}

// SSHClient dials up our target device.
func SSHClient(c SSHConfig) (*ssh.Client, error) {
	// Set up SSH auth method
	var sshAuth ssh.AuthMethod
	switch c.AuthMethod {
	case "KeyboardInteractive":
		sshAuth = ssh.KeyboardInteractive(c.sshInteractive)
		break
	case "PublicKeys":
		sshAuth = c.sshPublicKeys()
		break
	case "Password":
		sshAuth = ssh.Password(c.Pass)
		break
	// fail on unrecognized method
	default:
		return nil, fmt.Errorf("unrecognized SSH Authentication method: %s", c.AuthMethod)
	}

	// build SSH config
	config := &ssh.ClientConfig{
		User:    c.User,
		Auth:    []ssh.AuthMethod{sshAuth},
		Timeout: c.SSHTimeout * time.Second,
	}

	if c.AllowInsecure {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	} else {
		return nil, fmt.Errorf("host key checking not yet implimented")
		//config.HostKeyCallback = ssh.FixedHostKey(ssh.PublicKey(hostkey))
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
	if err != nil {
		log.Println("Unable to read SSH pivate key, reverting to password authentication")
		return ssh.Password(c.Pass)
	}
	privateKey, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Println("Unable to parse SSH private key, reverting to password authentication")
		return ssh.Password(c.Pass)
	}
	return ssh.PublicKeys(privateKey)
}
