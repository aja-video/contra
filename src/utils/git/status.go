package utils

import (
	"log"
)

//GitStatus reports the current working tree status
func GitStatus(g *Git) (bool, []string, error) {
	// Assume no changes
	gitChanged := false
	var updates []string
	s, err := g.Repo.Worktree()
	if err != nil {
		panic(err)
	}
	changes, err := s.Status()
	if ! changes.IsClean() {
		gitChanged = true
		for config := range changes {
			log.Printf("Config Changed: %s\n", config)
			updates = append(updates, config)
		}
	} else {
		log.Println("No config changes")
	}

	//log.Println(s.Status())
	//fmt.Printf("Current GIT status:\n%v", s.Status())
	return gitChanged, updates, err

}
