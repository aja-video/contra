package utils

import (
    "github.com/google/goexpect"
    "time"
    //	"regexp"
    "golang.org/x/crypto/ssh"
    "sort"
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

// BuildBatcher builds expect.Batcher from a send and receive map
func BuildBatcher(send map[int]string, receive map[int]string) *[]expect.Batcher {

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
        batch = append(batch, &expect.BExp{R:receive[k]})
        batch = append(batch, &expect.BSnd{S:send[k]})
    }

    // return
    return &batch
}