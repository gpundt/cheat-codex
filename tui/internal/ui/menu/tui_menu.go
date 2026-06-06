package tui_menu

import (
	Config "cheat-codex/internal/config"
	Styles "cheat-codex/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type MenuModel struct {
	ParentModel  *tea.MenuModel
	Choices		[]string
	Cursor int
	Width int
	Height int
}

func (model MenuModel) Init() tea.Cmd {
	return nil
}

func InitializeMenuModel() MenuModel {
	return MenuModel {
		ParentModel: nil,
		Choices: Config.GetActiveEmulators(),
		cursor: 0,
	}
}

func (model MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	selectedModel := "Menu"
	swtich msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "1":
			return model, tea.Quit
		
		case "up":
			if model.Cursor > 0 {
				model.Cursor--
			}
		case "down":
			if model.Cursor < len(model.Choices) - 1 {
				model.Cursor ++
			}
		case "enter", "space", "right":
			selectedModel = model.Choices[model.Cursor]
		}
		
	case tea.WindowSizeMsg:
		model.Width = msg.Width
		model.Height = msg.Height
		return model, nil 
	}

	var nextModel tea.Model

	nextModel = model
	return nextModel, nil
}

func (model MenuModel) View() string {
	title := Styles.Title.Render("Cheat Codex - ROM Emulator Hacking TUI")

	footer := Style.RenderFooter([][]string{
		{"↑↓", "navigate"},
		{"enter/space/→", "open"},
		{"ctrl+c/q", "quit"},
	})

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		footer,
	)
}