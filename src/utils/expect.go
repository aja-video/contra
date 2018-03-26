package utils

import (
	"fmt"
	"github.com/google/goexpect"
	"golang.org/x/crypto/ssh"
	"sort"
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

// BuildBatcher builds expect.Batcher from a send and receive map
func BuildBatcher(send map[int]string, receive map[int]string) []expect.Batcher {

	// Initialize stuff
	keys := make([]int, 0)
	batch := make([]expect.Batcher, 0)

	// Whichever is longer becomes our key list
	if len(receive) > len(send) {
		for k := range receive {
			keys = append(keys, k)
		}
	} else {
		for k := range send {
			keys = append(keys, k)
		}
	}

	// Sort, because golang
	sort.Ints(keys)

	// Set up batch, add send after every receive
	for _, k := range keys {
		batch = append(batch, &expect.BExp{R: receive[k]})
		batch = append(batch, &expect.BSnd{S: send[k]})
	}

	// return
	return batch
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
