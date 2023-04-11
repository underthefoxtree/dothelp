package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type buildToolsMainModel struct {
	options []string
	cursor  int
}

// Test
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
			if m.options[m.cursor] == "Build Tools" {
				return createBuildToolsMainModel(), nil
			} else {
				return m, tea.Quit
			}
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
			style = selectedItemStyle

			if choice == "Exit" {
				style = redItemStyle
			}
		}

		s += style.Render(fmt.Sprintf("%s%s", cursor, choice)) + "\n"
	}

	s += "\nPress q to quit."

	return s
}
