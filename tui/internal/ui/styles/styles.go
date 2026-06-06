package styles

import "github.com/charmbracelet/lipgloss"

var (
	// --- Color Palette ---------------------------
	colorPurple    = lipgloss.Color("99")
	colorDimPurple = lipgloss.Color("5")
	colorDimGray   = lipgloss.Color("236")
	colorWhite     = lipgloss.Color("255")


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

	
			// ── Footer ───────────────────────────────────────────────────
	FooterStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDimPurple).
			Background(colorDimGray).
			Width(ContentWidth).
			Padding(0, 2).
			MarginLeft(5)
)