package main

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type buildToolsMainModel struct {
	options []string
	cursor  int
}

func createBuildToolsMainModel() buildToolsMainModel {
	return buildToolsMainModel{
		options: []string{
			"Quick Build",
			"Release Build",
			"Complex Build",
			"Exit",
		},
		cursor: 0,
	}
}

func (m buildToolsMainModel) getOptionStyle() lipgloss.Style {
	style := selectedItemStyle

	if m.options[m.cursor] == "Exit" {
		style = redItemStyle
	}

	return style
}

func (m buildToolsMainModel) Init() tea.Cmd {
	return nil
}

func (m buildToolsMainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}

		case "enter", " ":
			switch m.cursor {
			case 0:
				return createSingleCommandModel(exec.Command("dotnet", "build"), "Build succeeded")
			case 1:
				return createSingleCommandModel(exec.Command("dotnet", "build", "-c", "Release"), "Release build succeeded.")
			case 2:
				return createExitModel("Complex builds aren't yet implemented.")
			}
			return createExitModel(
				fmt.Sprintf(
					"You selected: %s",
					m.getOptionStyle().Render(m.options[m.cursor])))
		}
	}

	return m, nil
}

func (m buildToolsMainModel) View() string {
	s := titleStyle.Render("DOTHELP") + "\n\n"

	for i, choice := range m.options {
		cursor := "  "
		style := itemStyle
		if i == m.cursor {
			cursor = "> "
			style = m.getOptionStyle()
		}

		s += cursor + style.Render(choice) + "\n"
	}

	s += helpStyle.Render("\nPress q to quit.")

	return s
}
