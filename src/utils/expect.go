package utils

import (
	"github.com/google/goexpect"
	"time"
	//	"regexp"
	"golang.org/x/crypto/ssh"
)

func GatherExpect(batcher *[]expect.Batcher, timeout time.Duration, ssh *ssh.Client) ([]expect.BatchRes, error) {

	// attach expect to our SSH connection
	ex, _, err := expect.SpawnSSH(ssh, timeout)

	if err != nil {
		panic(err)
	}

	// Gather data - batcher defined inside collector
	gather, err := ex.ExpectBatch(*batcher, timeout)

	if err != nil {
		panic(err)
	}

	// return config data
	return gather, err
}
