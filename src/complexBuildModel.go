package main

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type option struct {
	name  string
	value string
	otype int
}

type complexBuildModel struct {
	options []option
	cursor  int
	state   int
	list    tea.Model
	input   textinput.Model
	ch      chan string
}

func createComplexBuildModel() (tea.Model, tea.Cmd) {
	ch := make(chan string)

	return complexBuildModel{
		options: []option{
			{name: "Architecture", otype: 1},
			{name: "Force", otype: 0},
			{name: "Exit", otype: 2},
		},
		ch: ch,
	}, nil
}

func (m complexBuildModel) Init() tea.Cmd {
	return nil
}

func (m complexBuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case in := <-m.ch:
		switch m.options[m.cursor].otype {

		case 0:
			m.options[m.cursor].value = in
			m.state = 0
			return m, nil

		default:
			return m, tea.Quit
		}

	default:
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
					m.state = 1
					switch m.options[m.cursor].otype {
					case 0:
						m.list = listModel{
							options: []string{
								"Yes",
								"No",
							},
							getOptionStyle: func(option string) lipgloss.Style {
								return selectedItemStyle
							},
							selectOption: func(option string) (tea.Model, tea.Cmd) {
								go func() {
									m.ch <- option
								}()
								return m, tea.Every(time.Second, func(t time.Time) tea.Msg {
									return t
								})
							},
						}
					case 1:
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

		default:
			switch m.options[m.cursor].otype {

			case 0:
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd

			case 1:
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.String() {
					case "enter":
						m.options[m.cursor].value = m.input.Value()
						m.state = 0
						return m, nil
					}
				}

				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				return m, cmd

			default:
				return createExitModel("Exiting...")
			}
		}
	}
}

func (m complexBuildModel) View() string {
	switch m.state {
	case 0:
		s := titleStyle.Render("DOTHELP") + "\n\n"

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
				s += cursor + style.Render(choice.name, ":", v) + "\n"
			case 1:
				v := "Default"
				if choice.value != "" {
					v = choice.value
				}
				s += cursor + style.Render(choice.name, ":", v) + "\n"
			case 2:
				s += cursor + style.Render(choice.name) + "\n"
			}
		}

		s += helpStyle.Render("\nPress q to quit.")

		return s

	default:
		switch m.options[m.cursor].otype {
		case 0:
			return m.list.View()
		default:
			s := titleStyle.Render("DOTHELP") + "\n\n"

			s += "Set " + m.options[m.cursor].name + " (empty for default):\n"
			s += m.input.View()
			s += helpStyle.Render("\n\nPress q to quit")

			return s
		}
	}
}
