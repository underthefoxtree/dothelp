package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mainMenuModel struct {
	options []string
	cursor  int
}

func (m mainMenuModel) getOptionStyle() lipgloss.Style {
	style := selectedItemStyle

	if m.options[m.cursor] == "Exit" {
		style = redItemStyle
	}

	return style
}

func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return createExitModel("Exiting...")

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}

		case "enter", " ":
			switch m.options[m.cursor] {
			case "New Project":
				return createTemplateSelectModel(), nil
			case "Build Tools":
				return createBuildToolsMainModel(), nil
			default:
				return createExitModel(
					fmt.Sprintf(
						"You selected: %s",
						m.getOptionStyle().Render(m.options[m.cursor])))
			}
		}
	}

	return m, nil
}

func (m mainMenuModel) View() string {
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
