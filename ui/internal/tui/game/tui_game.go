package tui_game

import (
	Games "cheat-codex/internal/games"
	Memory "cheat-codex/internal/memory_map"
	Styles "cheat-codex/internal/tui/styles"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameModel struct {
	ParentModel    tea.Model
	SelectedGame   Games.Game
	OffsetEntryList []Memory.OffsetEntry
	Editing        bool
	EditInput      textinput.Model
	Cursor         int
	Width          int
	Height         int
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
	ti := textinput.New()
	ti.CharLimit = 128

	return GameModel{
		ParentModel:    parentModel,
		SelectedGame:   selectedGame,
		OffsetEntryList: selectedGame.Map.GenerateOffsetEntryList(),
		Editing:        false,
		EditInput:      ti,
		Cursor:         0,
		Width:          width,
		Height:         height,
	}
}

func (model GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	// --- Edit Mode --- Intercept all keys for the text input
	if model.Editing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// Commit the typed value
				switch model.OffsetEntryList[model.Cursor].Type {
				case "uint16":
					newEntry := model.OffsetEntryList[model.Cursor]
					newVal, err := strconv.ParseUint(
						model.EditInput.Value(), 10, 16,
					)
					if err != nil {
						break
					}
					newEntry.CurrentValue = int(newVal)
					model.SelectedGame.Map.UpdateOffsetEntryByOffset(
						newEntry.Offset.String(),
						newEntry,
					)
					model.OffsetEntryList = model.SelectedGame.Map.GenerateOffsetEntryList()
					model.Editing = false
					model.EditInput.Blur()
					return model, cmd
					
				case "uint8":
					newEntry := model.OffsetEntryList[model.Cursor]
					newVal, err := strconv.ParseUint(
						model.EditInput.Value(), 10, 8,
					)
					if err != nil {
						break
					}
					newEntry.CurrentValue = int(newVal)
					model.SelectedGame.Map.UpdateOffsetEntryByOffset(
						newEntry.Offset.String(),
						newEntry,
					)
					model.OffsetEntryList = model.SelectedGame.Map.GenerateOffsetEntryList()
					model.Editing = false
					model.EditInput.Blur()
					return model, cmd
				}
			case "esc":
				// Cancel — discard input
				model.Editing = false
				model.EditInput.Blur()
			default:
				model.EditInput, cmd = model.EditInput.Update(msg)
				// After update, refresh list of OffsetEntries
				model.OffsetEntryList = model.SelectedGame.Map.GenerateOffsetEntryList()
				return model, cmd
			}
		default:
			model.EditInput, cmd = model.EditInput.Update(msg)
			// After update, refresh list of OffsetEntries
			model.OffsetEntryList = model.SelectedGame.Map.GenerateOffsetEntryList()
			return model, cmd
		}

		// After update, refresh list of OffsetEntries
		model.OffsetEntryList = model.SelectedGame.Map.GenerateOffsetEntryList()
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
			if model.Cursor < len(model.OffsetEntryList)-1 {
				model.Cursor++
			}
			return model, nil
		case "enter", "space":
			switch model.OffsetEntryList[model.Cursor].Type {
			case "bool":
				// Calculate new boolean value
				newEntry := model.OffsetEntryList[model.Cursor]
				newEntry.CurrentValue ^= 1
				model.OffsetEntryList[model.Cursor].CurrentValue = newEntry.CurrentValue

				// Update the entry with matching offset
				model.SelectedGame.Map.UpdateOffsetEntryByOffset(
					newEntry.Offset.String(),
					newEntry,
				)

				// After update, refresh list of OffsetEntries
				model.OffsetEntryList = model.SelectedGame.Map.GenerateOffsetEntryList()
				return model, nil
			case "uint16", "uint8":
				// — enter edit mode ———————————————————
				model.EditInput.SetValue(
					fmt.Sprintf("%d", model.OffsetEntryList[model.Cursor].CurrentValue),
				)
				model.EditInput.Prompt = fmt.Sprintf(
					"%-10s:   ",
					Styles.SelectedItem.Render(
						strconv.Itoa(model.OffsetEntryList[model.Cursor].CurrentValue),
					),
				)
				model.EditInput.Focus()
				model.Editing = true
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
	var currentGroup = ""
	for _, offsetEntry := range model.OffsetEntryList {
		if currentGroup != offsetEntry.Group {
			currentGroup = offsetEntry.Group
			memoryMapGroup, _ := model.SelectedGame.Map.GetGroup(currentGroup)
			container = lipgloss.JoinVertical(
				lipgloss.Left,
				container,
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					Styles.GroupName.Render(memoryMapGroup.Name),
					Styles.GroupDescription.Render(memoryMapGroup.Description),
				),
			)
		}
		cursor := "  "
		if model.Cursor == rowNum {
			cursor = Styles.Cursor.Render("> ")
		}

		var entryValue = Styles.OffsetEntryValue.Render(
			fmt.Sprintf("%-10s", strconv.Itoa(offsetEntry.CurrentValue)),
		)
		if model.Editing && model.Cursor == rowNum {
			entryValue = Styles.SelectedItem.Render(
				fmt.Sprintf("%-10s", strconv.Itoa(offsetEntry.CurrentValue)),
			)
		}
		row := lipgloss.JoinHorizontal(
			lipgloss.Left,
			cursor,
			Styles.OffsetEntryLabel.Render(
				fmt.Sprintf("%-15s", offsetEntry.Label),
			),
			entryValue,
			Styles.OffsetEntryMisc.Render(
				fmt.Sprintf("%-6s", offsetEntry.Offset.String()),
			),
			Styles.OffsetEntryMisc.Render(
				fmt.Sprintf("%-10s", offsetEntry.Type),
			),
		)

		container = lipgloss.JoinVertical(
			lipgloss.Left,
			container,
			row,
		)
		rowNum++
		

	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		Styles.Container.Render(container),
		footer,
	)

}
