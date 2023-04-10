package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type programState int

const (
	quitting programState = iota
	mainMenu
	projectCreationMain
	buildTools
)

type model struct {
	state   programState
	options []string
	cursor  int
}

func initialModel() model {
	return model{
		state:   mainMenu,
		options: []string{"New Project", "Build Tools", "Exit"},
		cursor:  0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {

	// * Main Menu
	case mainMenu:
		switch msg := msg.(type) {

		case tea.KeyMsg:

			switch msg.String() {

			case "ctrl+c", "q":
				m.state = quitting
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
				if m.cursor != len(m.options)-1 {
					return m, tea.Quit
				} else {
					m.state = quitting
					return m, tea.Quit
				}
			}
		}

		return m, nil

	// * Quit
	case quitting:
		return m, tea.Quit

	// * Default
	default:
		return m, tea.Quit
	}
}

func (m model) View() string {
	switch m.state {

	// * Main Menu
	case mainMenu:
		s := "[---dothelp---]\n\n"

		for i, choice := range m.options {
			cursor := " "
			if i == m.cursor {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += "\nPress q to quit."

		return s

	// * Quit
	case quitting:
		return "Exiting dothelp...\n"

	// * Default
	default:
		return ""
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
