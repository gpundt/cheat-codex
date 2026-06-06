//go:build windows

package ipc

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

func listProcesses() ([]Process, error) {
	// tasklist /fo csv /nh gives "process.exe","1234","Console","1","10,000 K"
	out, err := exec.Command("tasklist", "/fo", "csv", "/nh").Output()
	if err != nil {
		return nil, err
	}

	var processes []Process
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Strip surrounding quotes and split on ","
		line = strings.Trim(line, `"`)
		parts := strings.Split(line, `","`)
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		pid, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		processes = append(processes, Process{PID: pid, Name: name})
	}

	log.Info().Str("func", "listProcesses").
		Msg(fmt.Sprintf("Found %d Windows processes", len(processes)))

	fmt.Printf("%v", processes)
	return processes, nil
}
