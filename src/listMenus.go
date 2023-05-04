package main

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func createBuildToolsMainModel() tea.Model {
	return listModel{
		options: []string{
			"Quick Build",
			"Release Build",
			"Complex Build",
			"Exit",
		},
		cursor: 0,
		getOptionStyle: func(option string) lipgloss.Style {
			if option == "Exit" {
				return redItemStyle
			} else {
				return selectedItemStyle
			}
		},
		selectOption: func(option string) (tea.Model, tea.Cmd) {
			switch option {
			case "Quick Build":
				return createSingleCommandModel(exec.Command("dotnet", "build"), "Build succeeded")
			case "Release Build":
				return createSingleCommandModel(exec.Command("dotnet", "build", "-c", "Release"), "Release build succeeded.")
			case "Complex Build":
				return createComplexBuildModel()
			default:
				return createExitModel(
					fmt.Sprintf(
						"You selected: %s",
						selectedItemStyle.Render(option)))
			}
		},
	}
}
