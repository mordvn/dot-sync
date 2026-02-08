package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorOK      = lipgloss.Color("2") // green
	ColorTODO    = lipgloss.Color("8") // gray
	ColorSyncing = lipgloss.Color("3") // yellow
	ColorError   = lipgloss.Color("1") // red
	ColorPrimary = lipgloss.Color("4") // blue
	ColorText    = lipgloss.Color("7") // whire
	ColorBG      = lipgloss.Color("0") // black

	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1).
			BorderForeground(lipgloss.Color("4"))

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("4")).
			MarginBottom(1)

	StyleItemOK = lipgloss.NewStyle().
			Foreground(ColorOK).
			Bold(true)

	StyleItemTODO = lipgloss.NewStyle().
			Foreground(ColorTODO)

	StyleItemSyncing = lipgloss.NewStyle().
				Foreground(ColorSyncing).
				Bold(true)

	StyleItemError = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	StyleFooter = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			MarginTop(1)

	StyleShortcut = lipgloss.NewStyle().
			Foreground(lipgloss.Color("4")).
			Bold(true)

	StyleHighlight = lipgloss.NewStyle().
			Background(lipgloss.Color("4")).
			Foreground(lipgloss.Color("0"))
)

func GetStatusStyle(status string) lipgloss.Style {
	switch status {
	case "OK":
		return StyleItemOK
	case "SYNCING":
		return StyleItemSyncing
	case "ERROR":
		return StyleItemError
	default:
		return StyleItemTODO
	}
}

func GetStatusIcon(status string) string {
	switch status {
	case "OK":
		return "✓"
	case "SYNCING":
		return "⟳"
	case "ERROR":
		return "✗"
	default:
		return "○"
	}
}