package utils

import "github.com/charmbracelet/lipgloss"

var ProviderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#B0197E")). // CSH Purple
	Padding(1, 2)

var ResultStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#E11C52")). // Hot Pink
	MarginLeft(2).
	Width(60)
