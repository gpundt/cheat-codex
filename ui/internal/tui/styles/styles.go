package styles

import "github.com/charmbracelet/lipgloss"

var (
	// --- Color Palette ---------------------------
	colorPurple    = lipgloss.Color("99")
	colorDimPurple = lipgloss.Color("5")
	colorGold      = lipgloss.Color("220")
	colorDimGray   = lipgloss.Color("236")
	colorWhite     = lipgloss.Color("255")
	colorDimWhite  = lipgloss.Color("250")
	colorGray      = lipgloss.Color("240")
	colorMuted     = lipgloss.Color("244")

	// ── Layout constants ─────────────────────────────────────────
	// 110 col terminal: 5 margin | 90 content | 5 margin
	ContentWidth = 100
	// Inner width = ContentWidth minus border (2) minus padding (4) = 134
	InnerWidth = 94

	// ── Title ───────────────────────────────────────────────────
	Title = lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPurple).
		Background(colorDimGray).
		Padding(0, 2).
		Foreground(colorWhite).
		Width(ContentWidth).
		MarginLeft(5)

	// ── Container ────────────────────────────────────────────────
	Container = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDimPurple).
			Align(lipgloss.Left).
			Width(ContentWidth).
			Padding(0, 2).
			Margin(0, 0, 1, 5) // bottom margin + left margin of 30

	ContainerHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGold).
			MarginBottom(1)

	// ── Footer ───────────────────────────────────────────────────
	FooterStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDimPurple).
			Background(colorDimGray).
			Width(ContentWidth).
			Padding(0, 2).
			MarginLeft(5)

	// ── Input ───────────────────────────────────────────────────
	Cursor = lipgloss.NewStyle().
		Foreground(colorGold).
		Bold(true)

	// ── List items ───────────────────────────────────────────────
	SelectedItem = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGold).
			PaddingLeft(2).
			PaddingRight(2)

	UnselectedItem = lipgloss.NewStyle().
			Foreground(colorDimWhite).
			PaddingLeft(2)
	
	Key = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPurple).
			Background(colorDimGray).
			Padding(0, 1)
	
	KeyDescription = lipgloss.NewStyle().
			Foreground(colorMuted).
			Background(colorDimGray).
			Padding(0, 1)

	// ── Helpers ───────────────────────────────────────────────────
	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPurple).
			Background(colorDimGray).
			Padding(0, 1)

	KeyDescStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Background(colorDimGray).
			Padding(0, 1)

	KeySepStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Background(colorDimGray)
)

func RenderFooter(binds [][]string) string {
	var parts []string
	sep := KeySepStyle.Render("·")

	for i, b := range binds {
		key := KeyStyle.Render(b[0])
		desc := KeyDescStyle.Render(b[1])
		parts = append(parts, key+desc)
		if i < len(binds)-1 {
			parts = append(parts, sep)
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Center, parts...)
	return FooterStyle.Render(row)
}
