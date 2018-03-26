package utils

import (
	"golang.org/x/crypto/ssh"
	"time"
)

// sshClient

func SshClient(user string, password string, host string) (*ssh.Client, error) {

	client, err := ssh.Dial("tcp", host, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	})

	if err != nil {
		panic(err)
	}

	return client, err

}
