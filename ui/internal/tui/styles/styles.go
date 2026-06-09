package styles

import "github.com/charmbracelet/lipgloss"

var (
	// --- Color Palette ---------------------------
	colorPurple    = lipgloss.Color("99")
	colorAqua      = lipgloss.Color("86")
	colorDimPurple = lipgloss.Color("5")
	colorGold      = lipgloss.Color("220")
	colorDimGray   = lipgloss.Color("236")
	colorWhite     = lipgloss.Color("255")
	colorDimWhite  = lipgloss.Color("250")
	colorGray      = lipgloss.Color("240")
	colorDimmerGray   = lipgloss.Color("241")
	colorMuted     = lipgloss.Color("244")

	// ── Layout constants ─────────────────────────────────────────


	// ── Title ───────────────────────────────────────────────────
	Title = lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPurple).
		Background(colorDimGray).
		Padding(0, 2).
		Foreground(colorWhite).
		// MarginLeft(10).
		MarginRight(10)

	// ── Container ────────────────────────────────────────────────
	Container = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDimPurple).
			Align(lipgloss.Left).
			Padding(0, 2).
			// MarginLeft(10).
			MarginRight(10)

	ContainerHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGold).
			Align(lipgloss.Left)

	// ── Logging ───────────────────────────────────────────────────
	InfoLogContainer = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Align(lipgloss.Left).
			Padding(0, 2).
			// MarginLeft(10).
			MarginRight(10)

	WarningLogContaner = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF7E00")).
			Align(lipgloss.Left).
			Padding(0, 0).
			// MarginLeft(10).
			MarginRight(10)
	
	ErrorLogContainer = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF0000")).
			Align(lipgloss.Left).
			Padding(0, 2).
			// MarginLeft(5).
			MarginRight(5)

	// ── Footer ───────────────────────────────────────────────────
	FooterStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDimPurple).
			Background(colorDimGray).
			Padding(0, 2).
			// MarginLeft(5).
			MarginRight(5)

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

	SelectedEditingItem = lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.RoundedBorder(), false, true, false, true).
			BorderForeground(colorPurple).
			Foreground(colorGold).
			PaddingLeft(0).
			PaddingRight(15)

	UnselectedItem = lipgloss.NewStyle().
			Foreground(colorDimWhite).
			PaddingLeft(2)

	// ── Memory Map Items ───────────────────────────────────────────────
	GroupName = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPurple).
			Background(colorDimGray).
			Padding(0, 1)
	
	GroupDescription = lipgloss.NewStyle().
			Foreground(colorMuted).
			Background(colorDimGray).
			Padding(0, 1)
	
	OffsetEntryLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAqua).
			PaddingRight(2).
			PaddingLeft(2)
	
	OffsetEntryMisc = lipgloss.NewStyle().
			Foreground(colorDimmerGray)
	
	OffsetEntryValue = lipgloss.NewStyle().
			Foreground(colorWhite).
			PaddingRight(2).
			PaddingLeft(2)
	
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

func RenderFooter(width int, binds [][]string) string {
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
	return FooterStyle.Width(width).MarginRight(5).Render(row)
}
