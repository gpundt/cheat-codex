package tui_game

import (
	Styles "cheat-codex/internal/tui/styles"
	Games "cheat-codex/internal/games"
	Memory "cheat-codex/internal/memory_map"
	"strconv"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

type GameModel struct {
	ParentModel tea.Model
	SelectedGame Games.Game
	AllOffsets []Memory.OffsetEntry
	Editing bool
	EditInput textinput.Model
	Cursor int
	Width int
	Height int
}

func (model GameModel) Init() tea.Cmd {
	return nil
}

func InitializeGameModel(
	parentModel tea.Model,
	selectedGame Games.Game,
	width,
	height int,
) GameModel {
	allOffsets := []Memory.OffsetEntry{}
	for _, group := range selectedGame.Map.Groups {
		for _, entry := range group.Offsets {
			allOffsets = append(allOffsets, entry)
		}
	}

	ti := textinput.New()
	ti.CharLimit = 128

	return GameModel{
		ParentModel: parentModel,
		SelectedGame: selectedGame,
		AllOffsets: allOffsets,
		Editing: false
		EditInput: ti,
		Cursor: 0,
		Width: width,
		Height: height,
	}
}

func (model GameModel) GetTotalRows() int {
	var totalRows = 0
	for _, group := range model.SelectedGame.Map.Groups {
		for _, _ = range group.Offsets {
			totalRows++
		}
	}
	return totalRows
}

func (model GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// --- Edit Mode --- Intercept all keys for the text input
	if model.Editing {
		return model, nil
	}
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return model, tea.Quit
		case "q", "left":
			return model.ParentModel, nil

		case "up":
			if model.Cursor > 0 {
				model.Cursor--
			}
		case "down":
			if model.Cursor < model.GetTotalRows() - 2 {
				model.Cursor++
			}
			return model, nil
		case "enter", "space":
			switch model.AllOffsets[model.Cursor + 1].Type {
			case "bool":	
				for groupNum, group := range model.SelectedGame.Map.Groups {
					for entryNum, entry := range group.Offsets {
						if entry.Offset == model.AllOffsets[model.Cursor + 1].Offset {
							model.SelectedGame.Map.Groups[groupNum].Offsets[entryNum].CurrentValue ^= 1
							
							// Now you gotta send a write
							return model, nil
						}
					}
				}
			case "uint16", "uint8":
				// — enter edit mode ———————————————————
				replacedKey := strings.ReplaceAll(row.ConfigName, "_", " ")
				titledKey := cases.Title(language.English).String(replacedKey)
				
				model.EditInput.Prompt = fmt.Sprintf(
					"%-42s:   ",
					Style.SelectedItemStyle.Render(titledKey),
				)
				model.EditInput.SetValue(fmt.Sprintf("%v", row.Value))
				model.EditInput.Focus()
				model.Editing = true
				return model, nil
			}

		}

	case tea.WindowSizeMsg:
		model.Width = msg.Width
		model.Height = msg.Height
		return model, nil
	}

	return model, nil
}

func (model GameModel) View() string {
	title := Styles.Title.Render(fmt.Sprintf(
		"%s Memory Modification",
		model.SelectedGame.Metadata.Name,
	))

	footer := Styles.RenderFooter([][]string{
		{"↑↓", "navigate"},
		{"enter/space", "Modify"},
		{"ctrl+c/esc", "quit"},
		{"←/q", "back"},
	})

	container := Styles.ContainerHeader.Render(
		"Modify individual memory addresses:",
	)
	container += lipgloss.JoinHorizontal(
		lipgloss.Left,
		fmt.Sprintf(
			"%-15s v%s",
			model.SelectedGame.Metadata.Name,
			model.SelectedGame.Metadata.Version,
		),
	)

	var rowNum = 0
	for _, group := range model.SelectedGame.Map.Groups {
		groupName := Styles.GroupName.Render(group.Name)
		groupDescription := Styles.GroupDescription.Render(group.Description)
		header := lipgloss.JoinHorizontal(
			lipgloss.Left,
			groupName,
			groupDescription,
		)
		container = lipgloss.JoinVertical(
				lipgloss.Left,
				container,
				header,
			)
		for _, offset := range group.Offsets {
			if offset.ReadOnly {
				continue
			}

			cursor := "  "
			if model.Cursor == rowNum {
				cursor = Styles.Cursor.Render("> ")
			}
			label := Styles.OffsetEntryLabel.Render(
				fmt.Sprintf("%-15s", offset.Label),
			)
			currentValue := Styles.OffsetEntryValue.Render(
				fmt.Sprintf("%-10s", strconv.Itoa(offset.CurrentValue)),
			)
			offsetString := Styles.OffsetEntryMisc.Render(
				fmt.Sprintf("%-6s", offset.Offset.String()),
			)
			valueType := Styles.OffsetEntryMisc.Render(
				fmt.Sprintf("%-10s", offset.Type),
			)
			row := lipgloss.JoinHorizontal(
				lipgloss.Left,
				cursor,
				label,
				currentValue,
				offsetString,
				valueType,
			)

			container = lipgloss.JoinVertical(
				lipgloss.Left,
				container,
				row,
			)
			rowNum++
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		Styles.Container.Render(container),
		footer,
	)

}