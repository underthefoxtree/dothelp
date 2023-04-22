package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type listModel struct {
	options        []string
	cursor         int
	getOptionStyle func(option string) lipgloss.Style
	selectOption   func(option string) (tea.Model, tea.Cmd)
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m.selectOption(m.options[m.cursor])
		}
	}

	return m, nil
}

func (m listModel) View() string {
	s := titleStyle.Render("DOTHELP") + "\n\n"

	for i, choice := range m.options {
		cursor := "  "
		style := itemStyle
		if i == m.cursor {
			cursor = "> "
			style = m.getOptionStyle(m.options[m.cursor])
		}

		s += cursor + style.Render(choice) + "\n"
	}

	s += helpStyle.Render("\nPress q to quit.")

	return s
}
