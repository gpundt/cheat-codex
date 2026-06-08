package tui_game

import (
	Styles "cheat-codex/internal/tui/styles"
	Games "cheat-codex/internal/games"
	"strconv"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameModel struct {
	ParentModel tea.Model
	SelectedGame Games.Game
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
	return GameModel{
		ParentModel: parentModel,
		SelectedGame: selectedGame,
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
			return model, nil
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
		groupName := Styles.Key.Render(group.Name)
		groupDescription := Styles.KeyDescription.Render(group.Description)
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
			label := Styles.Key.Render(offset.Label)
			valueType := Styles.Key.Render(offset.Type)
			currentValue := Styles.Key.Render(
				strconv.Itoa(offset.CurrentValue),
			)
			offset := Styles.Key.Render(offset.Offset.String())
			
			row := lipgloss.JoinHorizontal(
				lipgloss.Left,
				cursor,
				label,
				offset,
				valueType,
				currentValue,
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