package memory_map

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v3"
)

type HexOffset uint64

func (h *HexOffset) UmnarshalOffset(value *yaml.Node) error {
	raw := value.Value

	// Strip 0x prefix and parse hex
	if len(raw) > 2 && (raw[:2] == "0x" || raw[:2] == "0X") {
		parsed, err := strconv.ParseUint(raw[2:], 16, 64)
		if err != nil {
			return fmt.Errorf("Invalid hex offset %q: %w", raw, err)
		}
		*h = HexOffset(parsed)
		return nil
	}

	// Fall back to decimal
	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid offset %q: %w", raw, err)
	}
	*h = HexOffset(parsed)
	return nil
}

func (h HexOffset) String() string {
	return fmt.Sprintf("0x%x", uint64(h))
}

type OffsetEntry struct {
	Label        string    `yaml:"label"`
	Offset       HexOffset `yaml:"offset"`
	Type         string    `yaml:"type"`
	Minimum      int       `yaml:"min,omitempty"`
	Maximum      int       `yaml:"max,omitempty"`
	CurrentValue int
	ReadOnly     bool   `yaml:"readonly,omitempty"`
	Notes        string `yaml:"notes,omitempty"`
}

type Group struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Offsets     []OffsetEntry `yaml:"offsets"`
}

type MemoryMap struct {
	SchemaVersion int     `yaml:"schema_version"`
	Groups        []Group `yaml:"groups"`
}

