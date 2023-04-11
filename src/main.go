package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

// * Styles
var (
	titleStyle        = lipgloss.NewStyle().Bold(true).PaddingLeft(2)
	itemStyle         = lipgloss.NewStyle()
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.CompleteColor{
		TrueColor: "#8b8be1",
		ANSI256:   "62",
		ANSI:      "5",
	})
	redItemStyle = lipgloss.NewStyle().Foreground(lipgloss.CompleteColor{
		TrueColor: "#ff7f7f",
		ANSI256:   "203",
		ANSI:      "1",
	})

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.CompleteColor{
		TrueColor: "#3c3c3c",
		ANSI256:   "16",
		ANSI:      "0",
	})
)

// * Setup
func initialModel() mainMenuModel {
	return mainMenuModel{
		options: []string{"New Project", "Build Tools", "Exit"},
		cursor:  0,
	}
}

// * Main
func main() {
	program := tea.NewProgram(initialModel())

	if _, err := program.Run(); err != nil {
		fmt.Printf("An internal error occured: %v", err)
		os.Exit(1)
	}
}
