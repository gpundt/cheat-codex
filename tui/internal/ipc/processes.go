package ipc

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type Process struct {
	PID  int
	Name string
}

var EmualatorTargets = []string{"mgba", "melonds"}

func GetActiveEmulators() []Process {
	ActiveEmulators := []Process{}

	allProcesses, err := listProcesses()
	if err != nil {
		return ActiveEmulators
	}

	for _, p := range allProcesses {
		for _, target := range EmualatorTargets {
			if strings.Contains(
				strings.ToLower(p.Name),
				strings.ToLower(target),
			) {
				ActiveEmulators = append(ActiveEmulators, p)
				break
			}
		}
	}

	log.Info().Str("func", "GetActiveEmulators").
		Msg(fmt.Sprintf(
			"Found %d processes with matching substrings",
			len(ActiveEmulators),
		))
	return ActiveEmulators
}
