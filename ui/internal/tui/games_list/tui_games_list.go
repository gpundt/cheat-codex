package tui_games_list

import (
	Styles "cheat-codex/internal/tui/styles"
	Games "cheat-codex/internal/games"
	Game "cheat-codex/internal/tui/game"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GamesListModel struct {
	ParentModel tea.Model
	Emulator string
	Choices []Games.Game
	Cursor int
	Width int
	Height int
}

func (model GamesListModel) Init() tea.Cmd {
	return nil
}

func InitializeGamesListModel(
	parentModel tea.Model,
	emulator string,
	width,
	height int,
	) GamesListModel {
	return GamesListModel{
		ParentModel: parentModel,
		Emulator: emulator,
		Choices: Games.GetEmulatorGames(emulator),
		Cursor: 0,
		Width: width,
		Height: height,
	}
}

func (model GamesListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if model.Cursor < len(model.Choices)-1 {
				model.Cursor++
			}
			return model, nil
		case "enter", "space", "right":
			return Game.InitializeGameModel(
				model,
				model.Choices[model.Cursor],
				model.Width,
				model.Height,
			), nil
		}

	case tea.WindowSizeMsg:
		model.Width = msg.Width
		model.Height = msg.Height
		return model, nil
	}

	return model, nil
}

func (model GamesListModel) View() string {
	title := Styles.Title.Render(fmt.Sprintf(
		"Cheat Codex - %s Memory Modification",
		model.Emulator,
	))

	footer := Styles.RenderFooter([][]string{
		{"↑↓", "navigate"},
		{"enter/space/→", "select"},
		{"ctrl+c/esc", "quit"},
		{"←/q", "back"},
	})

	container := Styles.ContainerHeader.Render(fmt.Sprintf(
		"Select a game for the %s:",
		model.Emulator,
	))
	// container += fmt.Sprintf("%#v", model.Choices)

	if len(model.Choices) == 0 {
		container += fmt.Sprintf(
			"\nNo memory maps found for the %s emulator...",
			model.Emulator,
		)
	} else {
		for i, choice := range model.Choices {
			cursor := "  "
			gameName := choice.Metadata.Name
			mapSchema := strconv.Itoa(choice.Map.SchemaVersion)
			if model.Cursor == i {
				cursor = Styles.Cursor.Render("> ")
				gameName = Styles.SelectedItem.Render(
					fmt.Sprintf("%-18s", gameName),
				)
				mapSchema = Styles.SelectedItem.Render(
					fmt.Sprintf("Map Schema: %s", mapSchema),
				)
			} else {
				cursor = Styles.Cursor.Render("> ")
				gameName = Styles.UnselectedItem.Render(
					fmt.Sprintf("%-18s", gameName),
				)
				mapSchema = Styles.UnselectedItem.Render(
					fmt.Sprintf("Map Schema: %s", mapSchema),
				)
			}

			row := lipgloss.JoinHorizontal(
				lipgloss.Left,
				cursor,
				gameName,
				mapSchema,
			)
			container = lipgloss.JoinVertical(
				lipgloss.Left,
				container,
				row,
			)
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		Styles.Container.Render(container),
		footer,
	)
}