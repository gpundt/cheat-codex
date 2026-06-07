package ipc

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	Config "cheat-codex/internal/config"

	"github.com/rs/zerolog/log"
)

func callBinary(commandArguments string) string {
	args := strings.Fields(commandArguments)
	cmd := exec.Command(
		fmt.Sprintf(
			"./%s",
			Config.CodexCoreFilepath,
		),
		args...,
	)

	out, err := cmd.Output()

	if err != nil {
		// Check if binary exits with nonzero exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.Sys().(syscall.WaitStatus).ExitStatus()
			log.Err(
				errors.New("Binary exited with nonzero status code"),
			).Str("exit_code", strconv.Itoa(exitCode)).Str("func", "CallBinary").Msg("")
		} else {
			log.Err(err).Str("func", "CallBinary").Msg("")
		}
		return "error occurred"
	}

	return string(out)
}

func GetBaseAddress(pid int) string {
	log.Info().Str("pid", strconv.Itoa(pid)).Msg("Getting base address")

	commandArgs := fmt.Sprintf(
		"--action get-base-address --pid %d", pid,
	)
	result := strings.TrimSpace(callBinary(commandArgs))

	num, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		log.Err(err).Str("func", "GetBaseAddress").Msg("")
		return ""
	}

	hexStr := strconv.FormatInt(num, 16)

	return fmt.Sprintf("0x%s", hexStr)
}
