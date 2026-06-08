package config

var CodexFilepaths Filepaths

type Filepaths struct {
	CheatCodexEtc     string
	CheatCodexOpt     string
	CheatCodexBin     string
	CheatCodexSrc     string
	CoreBinary        string
	MapsDirectory     string
}

func InitializeFilepaths() Filepaths {
	CodexCoreFilepath := Filepaths{
		CheatCodexEtc: "/etc/cheat-codex/",
		CheatCodexOpt: "/opt/cheat-codex/",
		CheatCodexBin: "/opt/cheat-codex/bin/",
		CheatCodexSrc: "/opt/cheat-codex/src/",
		CoreBinary:    "/opt/cheat-codex/bin/" + CodexCoreFilepath
		MapsDirectory: "/opt/cheat-codex/maps/"
	}
}