package main

import (
	Config "cheat-codex/internal/config"
	Menu "cheat-codex/internal/tui/menu"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rs/zerolog/log"
)

func main() {
	Config.InitializeLogger()
	Config.InitializeFilepaths()

	program := tea.NewProgram(
		Menu.InitializeMenuModel(),
		tea.WithAltScreen(),
	)
	if _, err := program.Run(); err != nil {
		log.Fatal().Err(err).Str("func", "InitializeTUI").Msg("")
	}
}
