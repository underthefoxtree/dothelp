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
	greenItemStyle = lipgloss.NewStyle().Foreground(lipgloss.CompleteColor{
		TrueColor: "#3cb371",
		ANSI256:   "77",
		ANSI:      "2",
	})

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"})
)

// * Setup
func initialModel() tea.Model {
	return listModel{
		options: []string{"New Project", "Build Tools", "Exit"},
		cursor:  0,
		getOptionStyle: func(option string) lipgloss.Style {
			if option == "Exit" {
				return redItemStyle
			} else {
				return selectedItemStyle
			}
		},
		selectOption: func(option string) (tea.Model, tea.Cmd) {
			switch option {
			case "New Project":
				return createTemplateSelectModel(), nil
			case "Build Tools":
				return createBuildToolsMainModel(), nil
			default:
				return createExitModel(
					fmt.Sprintf(
						"You selected: %s",
						selectedItemStyle.Render(option)))
			}
		},
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
