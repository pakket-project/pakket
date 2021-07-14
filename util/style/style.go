package style

import "github.com/charmbracelet/lipgloss"

var (
	Pkg = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12"))

	Repo = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12"))

	// Version = lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("198"))

	Link = lipgloss.NewStyle().
		Underline(true).
		Foreground(lipgloss.Color("#4aaada"))

	Error = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9"))

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))
)
