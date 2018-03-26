package utils

import (
	"fmt"
	"github.com/google/goexpect"
	"golang.org/x/crypto/ssh"
	"time"
)

// GatherExpect initializes expect with our SSH connection and gathers results of the batch.
func GatherExpect(batcher []expect.Batcher, timeout time.Duration, ssh *ssh.Client) ([]expect.BatchRes, error) {

	// attach expect to our SSH connection
	ex, _, err := expect.SpawnSSH(ssh, timeout)

	if err != nil {
		panic(err)
	}

	// Gather data - batcher defined inside collector
	gather, err := ex.ExpectBatch(batcher, timeout)

	if err != nil {
		panic(err)
	}

	// return config data
	return gather, err
}

// SimpleBatcher implements a straight forward send/receive pattern.
func SimpleBatcher(definition [][]string) ([]expect.Batcher, error) {
	batchSlice := make([]expect.Batcher, 0)

	for _, set := range definition {
		// Make sure we have either 1 or 2 values.
		if len(set) < 1 || len(set) > 2 {
			return nil, fmt.Errorf("definition expects 1 or 2 values, but found: %d", len(set))
		}

		// Simple batcher always expects the first value to be a receive.
		batchSlice = append(batchSlice, &expect.BExp{R: set[0]})
		if len(set) == 2 {
			// If a second value is provided, we plan to send it.
			batchSlice = append(batchSlice, &expect.BSnd{S: set[1]})
		}
	}

	return batchSlice, nil
}
