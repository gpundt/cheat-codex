package tui_game

import (
	Config "cheat-codex/internal/config"
	Games "cheat-codex/internal/games"
	Memory "cheat-codex/internal/memory_map"
	Styles "cheat-codex/internal/tui/styles"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameModel struct {
	ParentModel  tea.Model
	SelectedGame Games.Game
	TableRows    []Memory.TableRow
	Viewport     viewport.Model
	Editing      bool
	EditInput    textinput.Model
	LogMessage   *Config.LogStruct
	Cursor       int
	Width        int
	Height       int
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

	model := GameModel{
		ParentModel:  parentModel,
		SelectedGame: selectedGame,
		TableRows:    selectedGame.Map.GetTableRows(),
		Editing:      false,
		EditInput:    ti,
		LogMessage:   nil,
		Cursor:       0,
		Width:        width,
		Height:       height,
	}

	vp := viewport.New(width-10, 20)
	vp.SetContent(model.generateTableContent())
	model.Viewport = vp
	return model
}

func (model GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// --- Edit Mode --- Intercept all keys for the text input
	if model.Editing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return model, tea.Quit
			case "enter":
				// Commit the typed value
				switch model.TableRows[model.Cursor].Type {
				case "uint16", "uint8":
					var err error = nil
					if model.TableRows[model.Cursor].Type == "uint16" {
						_, err = strconv.ParseUint(
							model.EditInput.Value(), 10, 16,
						)
					} else {
						_, err = strconv.ParseUint(
							model.EditInput.Value(), 10, 8,
						)
					}
					if err != nil {
						model.LogMessage = &Config.LogStruct{
							Severity: "ERROR",
							Message:  err.Error(),
						}
						return model, nil
					}

					model.TableRows[model.Cursor].CurrentValue = model.EditInput.Value()
					if err := model.SelectedGame.Map.UpdateMapFromTableRows(
						model.TableRows[model.Cursor],
					); err != nil {
						model.LogMessage = &Config.LogStruct{
							Severity: "ERROR",
							Message:  err.Error(),
						}
					}

					model.Editing = false
					model.EditInput.Blur()

					model.LogMessage = &Config.LogStruct{
						Severity: "INFO",
						Message: fmt.Sprintf(
							"Set %s to %s",
							model.TableRows[model.Cursor].Label,
							model.TableRows[model.Cursor].CurrentValue,
						),
					}
					model.Viewport.SetContent(model.generateTableContent())
					return model, cmd
				}
			case "esc":
				// Cancel — discard input
				model.Editing = false
				model.EditInput.Blur()
				return model, nil

			default:
				model.EditInput, cmd = model.EditInput.Update(msg)
				model.TableRows[model.Cursor].CurrentValue = model.EditInput.Value()
				if err := model.SelectedGame.Map.UpdateMapFromTableRows(
					model.TableRows[model.Cursor],
				); err != nil {
					model.LogMessage = &Config.LogStruct{
						Severity: "ERROR",
						Message:  err.Error(),
					}
				}
				model.Viewport.SetContent(model.generateTableContent())
				return model, cmd
			}
		default:
			model.EditInput, cmd = model.EditInput.Update(msg)
			return model, cmd
		}

		// After update, refresh list of OffsetEntries
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
				model.Viewport.SetContent(model.generateTableContent())
			}


		case "down":
			if model.Cursor < len(model.TableRows)-1 {
				model.Cursor++
				model.Viewport.SetContent(model.generateTableContent())
			}

		case "enter", "space":
			switch model.TableRows[model.Cursor].Type {
			case "bool":
				// Calculate new boolean value and update
				valueInt, _ := strconv.Atoi(model.TableRows[model.Cursor].CurrentValue)
				newValue := strconv.Itoa(valueInt ^ 1)
				model.TableRows[model.Cursor].CurrentValue = newValue

				// Update the entry with matching offset
				if err := model.SelectedGame.Map.UpdateMapFromTableRows(
					model.TableRows[model.Cursor],
				); err != nil {
					model.LogMessage = &Config.LogStruct{
						Severity: "ERROR",
						Message:  err.Error(),
					}
				}

				model.LogMessage = &Config.LogStruct{
					Severity: "INFO",
					Message: fmt.Sprintf(
						"Set %s to %s",
						model.TableRows[model.Cursor].Label,
						model.TableRows[model.Cursor].CurrentValue,
					),
				}
				model.Viewport.SetContent(model.generateTableContent())

			case "uint16", "uint8":
				// — enter edit mode ———————————————————
				model.EditInput.SetValue(
					model.TableRows[model.Cursor].CurrentValue,
				)
				model.EditInput.Prompt = fmt.Sprintf(
					"%-10s:   ",
					Styles.SelectedItem.Render(
						model.TableRows[model.Cursor].CurrentValue,
					),
				)

				model.EditInput.Focus()
				model.Editing = true
				model.Viewport.SetContent(model.generateTableContent())
			}
		}

	case tea.WindowSizeMsg:
		model.Viewport.Width = msg.Width - 10
		model.Viewport.Height = msg.Height - 30
	}

	model.Viewport, cmd = model.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	return model, tea.Batch(cmds...)
}

func (model GameModel) View() string {
	title := Styles.Title.Width(model.Width - 10).Render(fmt.Sprintf(
		"%s Memory Modification",
		model.SelectedGame.Metadata.Name,
	))

	footer := Styles.RenderFooter(model.Width-10, [][]string{
		{"↑↓", "navigate"},
		{"enter/space", "Modify"},
		{"ctrl+c/esc", "quit"},
		{"←/q", "back"},
	})

	container := model.Viewport.View()

	var logContainer = Styles.InfoLogContainer.Width(model.Width - 10).Render("")
	if model.LogMessage != nil {
		switch model.LogMessage.Severity {
		case "INFO":
			logContainer = Styles.InfoLogContainer.Width(model.Width - 10).Render(
				model.LogMessage.Message,
			)
		case "WARN":
			logContainer = Styles.WarningLogContaner.Width(model.Width - 10).Render(
				model.LogMessage.Message,
			)
		case "ERROR":
			logContainer = Styles.ErrorLogContainer.Width(model.Width - 10).Render(
				model.LogMessage.Message,
			)
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		container,
		logContainer,
		footer,
	)
}

func (model GameModel) generateTableContent() string {
	container := Styles.ContainerHeader.Width(model.Width - 14).Render(
		"Modify individual memory addresses:",
	)
	container += fmt.Sprintf(
		"\n%-15s v%s\n",
		model.SelectedGame.Metadata.Name,
		model.SelectedGame.Metadata.Version,
	)

	var rowNum = 0
	var currentGroup = ""
	for _, tableRow := range model.TableRows {
		if currentGroup != tableRow.Group {
			currentGroup = tableRow.Group
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
			fmt.Sprintf("%-14s", tableRow.CurrentValue),
		)
		if model.Editing && model.Cursor == rowNum {
			entryValue = Styles.SelectedEditingItem.Render(
				tableRow.CurrentValue,
			)
		}

		renderedRow := lipgloss.JoinHorizontal(
			lipgloss.Left,
			cursor,
			Styles.OffsetEntryLabel.Render(fmt.Sprintf("%-15s", tableRow.Label)),
			entryValue,
			Styles.OffsetEntryMisc.Render(fmt.Sprintf("%-14s", tableRow.Offset)),
			Styles.OffsetEntryMisc.Render(fmt.Sprintf("%-10s", tableRow.Type)),
		)

		container = lipgloss.JoinVertical(
			lipgloss.Left,
			container,
			renderedRow,
		)
		rowNum++
	}

	return Styles.Container.Width(model.Width - 12).Render(container)
}
