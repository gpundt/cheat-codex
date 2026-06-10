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
	Group        string
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

// Helper function to swap na offset entry with a matching offset with a new OffsetEntry
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

type TableRow struct {
	Group        string
	Label        string
	CurrentValue string
	Offset       string
	Type         string
}

func (mm MemoryMap) GetTableRows() []TableRow {
	rows := []TableRow{}
	for _, group := range mm.Groups {
		for _, entry := range group.Offsets {
			if entry.ReadOnly {
				continue
			}

			rows = append(rows, TableRow{
				Group:        group.Name,
				Label:        entry.Label,
				CurrentValue: strconv.Itoa(entry.CurrentValue),
				Offset:       entry.Offset.String(),
				Type:         entry.Type,
			})
		}
	}

	return rows
}

func (mm MemoryMap) UpdateMapFromTableRows(row TableRow) error {
	for groupIndex := range mm.Groups {
		for entryIndex := range mm.Groups[groupIndex].Offsets {
			entry := &mm.Groups[groupIndex].Offsets[entryIndex]
			if entry.ReadOnly {
				continue
			}

			if entry.Offset.String() == row.Offset {
				var err error
				val, err := strconv.Atoi(row.CurrentValue)
				if err != nil {
					return err
				}
				entry.CurrentValue = val
				return nil
			}
		}
	}

	return fmt.Errorf("No Map entry with offset '%s' found", row.Offset)
}

func GetOffsetEntriesByGroup(
	entryList []OffsetEntry,
	groupName string,
) []OffsetEntry {
	entriesOfGroup := []OffsetEntry{}
	for _, entry := range entryList {
		if entry.Group == groupName {
			entriesOfGroup = append(entriesOfGroup, entry)
		}
	}

	return entriesOfGroup
}
