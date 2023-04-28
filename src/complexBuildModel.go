package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type option struct {
	name  string
	value string
	otype int
}

type complexBuildModel struct {
	options    []option
	cursor     int
	optionBool bool
	state      int
	input      textinput.Model
}

func createComplexBuildModel() (tea.Model, tea.Cmd) {
	return complexBuildModel{
		options: []option{
			{name: "Architecture", otype: 1},
			{name: "Force", otype: 0},
			{name: "Exit", otype: 2},
		},
	}, nil
}

func (m complexBuildModel) Init() tea.Cmd {
	return nil
}

func (m complexBuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case 0:
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

			case "enter":
				switch m.options[m.cursor].otype {
				case 0:
					m.state = 1
					m.optionBool = true

				case 1:
					m.state = 2
					m.input = textinput.New()
					m.input.Placeholder = "Default"
					m.input.Focus()
					return m, textinput.Blink

				default:
					return createExitModel("Exiting...")
				}
			}
		}

		return m, nil

	case 1:
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k", "down", "j":
				m.optionBool = !m.optionBool
			case "enter":
				m.state = 0
				if m.optionBool {
					m.options[m.cursor].value = "Yes"
				} else {
					m.options[m.cursor].value = "No"
				}
			}
		}

		return m, nil

	case 2:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return createExitModel("Exiting...")
			case "enter":
				m.state = 0
				m.options[m.cursor].value = m.input.Value()
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd

	default:
		return createExitModel("Exit")
	}
}

func (m complexBuildModel) View() string {
	s := titleStyle.Render("DOTHELP") + "\n\n"

	switch m.state {
	case 0:
		for i, choice := range m.options {
			cursor := "  "
			style := itemStyle
			if i == m.cursor {
				cursor = "> "
				style = selectedItemStyle

				if choice.otype == 2 {
					style = redItemStyle
				}
			}

			switch choice.otype {
			case 0:
				v := "No"
				if choice.value != "" {
					v = choice.value
				}
				s += cursor + style.Render(choice.name+":", v) + "\n"
			case 1:
				v := "Default"
				if choice.value != "" {
					v = choice.value
				}
				s += cursor + style.Render(choice.name+":", v) + "\n"
			case 2:
				s += cursor + style.Render(choice.name) + "\n"
			}
		}

	case 1:
		if m.optionBool {
			s += "> " + selectedItemStyle.Render("Yes") + "\n  " + itemStyle.Render("No")
		} else {
			s += "  " + itemStyle.Render("Yes") + "\n> " + selectedItemStyle.Render("No")
		}

	case 2:
		s += m.input.View()
	}

	s += helpStyle.Render("\nPress q to quit.")
	return s
}
