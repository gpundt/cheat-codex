package tui_process

import (
	IPC "cheat-codex/internal/ipc"
	MemoryMap "cheat-codex/internal/memory_map"
	Styles "cheat-codex/internal/ui/styles"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EmulatorModel struct {
	ParentModel *tea.Model
	Emulator    IPC.EmulatorProcess
	Map         MemoryMap.MemoryMap
	Cursor      int
	Width       int
	Height      int
}

func (model EmulatorModel) Init() tea.Cmd {
	return nil
}

func InitializeEmulatorModel(
	parentModel tea.Model,
	emulator IPC.EmulatorProcess,
	width,
	height int,
) EmulatorModel {
	mm, err := MemoryMap.InitializeMemoryMap()
	return EmulatorModel{
		ParentModel: &parentModel,
		Emulator:    emulator,
		Map:         MemoryMap.InitializeMemoryMap(),
		Cursor:      0,
		Width:       width,
		Height:      height,
	}
}

func (model EmulatorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return model, tea.Quit
		case "q", "left":
			return *model.ParentModel, nil

		case "up":
			if model.Cursor > 0 {
				model.Cursor--
			}
		case "down":
			// if model.Cursor < len(model.Choices)-1 {
			// 	model.Cursor++
			// }
			return model, nil
		case "enter", "space", "right":
			return model, nil
		}

	case tea.WindowSizeMsg:
		model.Width = msg.Width
		model.Height = msg.Height
		return model, nil
	}

	return model, nil
}

func (model EmulatorModel) View() string {
	title := Styles.Title.Render("Cheat Codex - Memory Modification")

	footer := Styles.RenderFooter([][]string{
		{"↑↓", "navigate"},
		{"enter/space", "modify"},
		{"ctrl+c/esc", "quit"},
		{"←/q", "back"},
	})

	container := Styles.ContainerHeader.Render(fmt.Sprintf(
		"Process: %-15sPID: %-8sEmulator Base Address: %-15sMemory Region Base Address: %-15s",
		model.Emulator.Name,
		strconv.Itoa(model.Emulator.PID),
		model.Emulator.EmulatorBaseAddress,
		model.Emulator.RegionBaseAddress,
	))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		Styles.Container.Render(container),
		footer,
	)
}
