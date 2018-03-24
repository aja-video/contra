package utils

import (
	"github.com/google/goexpect"
	"regexp"
	"time"
)


func Expect(cmd string, timeout time.Duration, target string) (string, []string, error){
	// spawn our command
	ex, _, err := expect.Spawn(cmd, timeout)
	if err != nil {
		panic(err)
	}
	// parse output
	out, match, err := ex.Expect(regexp.MustCompile(target), timeout)

	if err != nil {
		panic(err)
	}

	return out, match, err
}