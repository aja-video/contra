package utils

import (
	"gopkg.in/src-d/go-git.v4/plumbing"
)

//GitStatus reports the current working tree status
func GitStatus(g *Git) (*plumbing.Reference, error) {

	return g.Repo.Head() // TODO this may be adequate, but "changed" is something I was able to pull testing before

}
