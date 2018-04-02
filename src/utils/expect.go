package utils

import (
	"fmt"
	"github.com/google/goexpect"
	"golang.org/x/crypto/ssh"
	"regexp"
	"strconv"
	"time"
)

// GatherExpect initializes expect with our SSH connection and gathers results of the batch.
func GatherExpect(batcher []expect.Batcher, timeout time.Duration, ssh *ssh.Client) ([]expect.BatchRes, error) {

	// attach expect to our SSH connection
	ex, _, err := expect.SpawnSSH(ssh, timeout)
	// Use the below for verbose debugging output. Should only be used for debugging since it may
	// dump passwords to the screen depending on the collector used.
	//ex, _, err := expect.SpawnSSH(ssh, timeout, expect.Verbose(true), expect.VerboseWriter(os.Stdout))

	if err != nil {
		return nil, err
	}

	// Gather data - batcher defined inside collector
	gather, err := ex.ExpectBatch(batcher, timeout)

	if err != nil {
		return nil, err
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

// VariableBatcher uses send and receive, with the possibility to skip steps.
// TODO: Do we need https://godoc.org/google.golang.org/grpc/codes ?
func VariableBatcher(definition [][]string) ([]expect.Batcher, error) {
	caseSlice := make([]expect.Caser, 0)

	for i, set := range definition {
		// Make sure we have either 1 or 2 values.
		if len(set) < 1 || len(set) > 2 {
			return nil, fmt.Errorf("definition expects 1 or 2 values, but found: %d", len(set))
		}

		// Always have a receive to match.
		batchCase := &expect.Case{R: regexp.MustCompile(set[0])}
		if len(set) == 2 {
			batchCase.S = set[1]
			// If this matches more than Rt times, fail with the number of the case.
			batchCase.T = expect.Continue(expect.NewStatus(7, "failed case: "+strconv.Itoa(i)))
			// For this VariableBatcher, we will only match once.
			batchCase.Rt = 1
		} else {
			// The success case.
			batchCase.T = expect.OK()
		}
		caseSlice = append(caseSlice, batchCase)
	}

	// Put all the cases in a BCas command, and put the Caser in the command.
	batchSlice := []expect.Batcher{&expect.BCas{C: caseSlice}}

	return batchSlice, nil
}

// ComplexBatcher will implement more complex switch cases.
//func ComplexBatcher() {}
