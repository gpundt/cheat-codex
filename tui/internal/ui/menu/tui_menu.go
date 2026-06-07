package tui_menu

import (
	IPC "cheat-codex/internal/ipc"
	Process "cheat-codex/internal/ui/process"
	Styles "cheat-codex/internal/ui/styles"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuModel struct {
	ParentModel *tea.Model
	Choices     []IPC.Process
	Cursor      int
	Width       int
	Height      int
}

func (model MenuModel) Init() tea.Cmd {
	return nil
}

func InitializeMenuModel() MenuModel {
	return MenuModel{
		ParentModel: nil,
		Choices:     IPC.GetActiveEmulators(),
		Cursor:      0,
	}
}

func (model MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return model, tea.Quit

		case "up":
			if model.Cursor > 0 {
				model.Cursor--
			}
		case "down":
			if model.Cursor < len(model.Choices)-1 {
				model.Cursor++
			}
		case "enter", "space", "right":
			emulatorProcess := IPC.EmulatorProcess{
				Name:                model.Choices[model.Cursor].Name,
				PID:                 model.Choices[model.Cursor].PID,
				EmulatorBaseAddress: IPC.GetEmulatorBaseAddress(model.Choices[model.Cursor].PID),
				RegionBaseAddress:   "0x00",
			}
			return Process.InitializeEmulatorModel(
				model,
				emulatorProcess,
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

func (model MenuModel) View() string {
	title := Styles.Title.Render("Cheat Codex - ROM Emulator Hacking TUI")

	footer := Styles.RenderFooter([][]string{
		{"↑↓", "navigate"},
		{"enter/space/→", "open"},
		{"ctrl+c/q/esc", "quit"},
	})

	container := Styles.ContainerHeader.Render(
		"Select an emulator process to attach to:\n",
	)

	if len(model.Choices) == 0 {
		container += "\nNo emulator processes detected..."
	} else {
		for i, choice := range model.Choices {
			cursor := "  "
			name := choice.Name
			pid := strconv.Itoa(choice.PID)
			if model.Cursor == i {
				cursor = Styles.Cursor.Render("> ")
				name = Styles.SelectedItem.Render(
					fmt.Sprintf("%-18s", name),
				)
				pid = Styles.SelectedItem.Render(
					fmt.Sprintf("PID:%s", pid),
				)
			} else {
				name = Styles.UnselectedItem.Render(
					fmt.Sprintf("%-20s", name),
				)
				pid = Styles.UnselectedItem.Render(
					fmt.Sprintf("PID:%s", pid),
				)
			}

			row := lipgloss.JoinHorizontal(
				lipgloss.Left,
				cursor,
				name,
				pid,
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
