package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

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
