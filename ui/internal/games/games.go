package games

import (
	"os"
	"strings"
	"path/filepath"

	MemoryMap "cheat-codex/internal/memory_map"

	"github.com/rs/zerolog/log"
	"go.yaml.in/yaml/v3"
)

type Game struct {
	Metadata struct {
		Name     string `yaml:"name"`
		Version  string `yaml:"version,omitempty"`
		Platform string `yaml:"platform"`
		Emulator string `yaml:"emulator"`
		Process  string `yaml:"process"`
		Region   string `yaml:"region"`
	} `yaml:"meta"`
	Map MemoryMap.MemoryMap
}

func InitializeGameStruct(mapFilepath string) Game {
	data, err := os.ReadFile(mapFilepath)
	if err != nil {
		log.Fatal().Err(err).Str("func", "InitializeGameStruct").
			Msg("Failed to read file")
	}
	
	var mm MemoryMap.MemoryMap
	if err := yaml.Unmarshal(data, &mm); err != nil {
		log.Fatal().Err(err).Str("func", "InitializeGameStruct").
			Msg("Failed to unmarshal MemoryMap")
	}
	
	var game Game
	if err := yaml.Unmarshal(data, &game); err != nil {
		log.Fatal().Err(err).Str("func", "InitializeGameStruct").
			Msg("Failed to unmarshal GameMetadata")
	}

	game.Map = mm

	return game
}

func GetEmulatorGames(emulatorName string) []Game {
	games := []Game{}

	mapsDir := "/opt/cheat-codex/maps"

	// Read the directory contents
	entries, err := os.ReadDir(mapsDir)
	if err != nil {
		log.Fatal().Err(err).Str("func", "GetEmulatorGames").
			Msg("Failed to get directory contents")
	}

	for _, entry := range entries {
		if strings.Contains(entry.Name(), emulatorName) {
			games = append(
				games,
				InitializeGameStruct(filepath.Join(mapsDir, entry.Name())),
			)
		}
	}

	return games
}