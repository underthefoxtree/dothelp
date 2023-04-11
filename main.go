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
)

// * Models
type mainMenuModel struct {
	options []string
	cursor  int
}

type buildToolsMainModel struct {
	options []string
	cursor  int
}

type exitModel struct {
	exitMessage string
}

// * Exit Model
func createExitModel(msg string) exitModel {
	return exitModel{
		exitMessage: msg,
	}
}

func (m exitModel) Init() tea.Cmd {
	return tea.Quit
}

func (m exitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m exitModel) View() string {
	return m.exitMessage + "\n"
}

// * Main Menu
func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return createExitModel("Exiting..."), tea.Quit

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

func (m mainMenuModel) View() string {
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

// * Build Tools Main
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
