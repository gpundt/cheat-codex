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

// Helper function to get individual group with a matching name
func (mm MemoryMap) GetGroup(name string) (*Group, error) {
	for _, group := range mm.Groups {
		if group.Name == name {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("No group with name '%s' found...", name)
}

// Helper function to get an individual OffsetEntry with a matching label
func (mm MemoryMap) GetOffsetEntry(
	label string,
) (*OffsetEntry, error) {
	for _, group := range mm.Groups {
		for _, entry := range group.Offsets {
			if entry.Label == label {
				return &entry, nil
			}
		}
	}

	return nil, fmt.Errorf("No offset entry with label '%s' found...", label)
}

// Helper function to get a list of all OffsetEntries from the memory map
func (mm MemoryMap) GetAllOffsetEntries() (int, []OffsetEntry) {
	entries := []OffsetEntry{}
	for _, group := range mm.Groups {
		for _, entry := range group.Offsets {
			entries = append(entries, entry)
		}
	}

	return len(entries), entries
}

// Helper function to swap an individual OffsetEntry of a specific label with another OffsetEntry
func (mm MemoryMap) UpdateOffsetEntryByLabel(
	label string,
	newEntry OffsetEntry,
) error {
	for groupIndex, group := range mm.Groups {
		for entryIndex, entry := range group.Offsets {
			if entry.Label == label {
				mm.Groups[groupIndex].Offsets[entryIndex] = newEntry
				return nil
			}
		}
	}

	return fmt.Errorf("No offset entry with label '%s' found...", label)
}

// Helper function to swap
func (mm MemoryMap) UpdateOffsetEntryByOffset(
	offset string,
	newEntry OffsetEntry,
) error {
	for groupIndex, group := range mm.Groups {
		for entryIndex, entry := range group.Offsets {
			if entry.Offset.String() == offset {
				mm.Groups[groupIndex].Offsets[entryIndex] = newEntry
				return nil
			}
		}
	}

	return fmt.Errorf("No offset entry with offset '%s' found...", offset)
}
