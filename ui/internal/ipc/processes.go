package ipc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type EmulatorProcess struct {
	ProcessName         string
	EmulatorName		string
	PID                 int
	EmulatorBaseAddress string
	RegionBaseAddress   string
}

type Process struct {
	PID  int
	Name string
}

var EmualatorTargets = []string{"mgba", "melonds"}

func GetActiveEmulators() []EmulatorProcess {
	ActiveEmulators := []EmulatorProcess{}

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
				emulator := EmulatorProcess{
					ProcessName: p.Name,
					EmulatorName: target,
					PID: p.PID,
				}
				ActiveEmulators = append(ActiveEmulators, emulator)
				break
			}
		}
	}

	log.Info().Str("func", "GetActiveEmulators").
		Str("matching_processes", strconv.Itoa(len(ActiveEmulators))).
		Msg(fmt.Sprintf(
			"%v", ActiveEmulators,
		))
	return ActiveEmulators
}
