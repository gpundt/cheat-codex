//go:build linux

package ipc

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// Function to collect all running processes on a linux host
func listProcesses() ([]Process, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	var processes []Process
	for _, entry := range entries {
		// Isolate PIDs
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || !entry.IsDir() {
			continue
		}

		// Get process name and trim
		commPath := filepath.Join("/proc", entry.Name(), "comm")
		data, err := os.ReadFile(commPath)
		if err != nil {
			// Skip processes that end between ReadDir and now
			continue
		}

		name := strings.TrimSpace(string(data))
		processes = append(processes, Process{PID: pid, Name: name})
	}

	log.Info().Str("func", "listProcesses").
		Msg(fmt.Sprintf("Found %d Linux processes", len(processes)))
	return processes, nil
}
