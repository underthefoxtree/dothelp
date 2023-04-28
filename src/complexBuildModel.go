package main

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type option struct {
	name  string
	value string
	flag  string
	otype int
}

type complexBuildModel struct {
	options    []option
	cursor     int
	optionBool bool
	state      int
	input      textinput.Model
	paginator  paginator.Model
}

func createComplexBuildModel() (tea.Model, tea.Cmd) {
	options := []option{
		{name: "Architecture", otype: 1, flag: "-a"},
		{name: "Configuration", otype: 1, flag: "-c"},
		{name: "Framework", otype: 1, flag: "-f"},
		{name: "Force", otype: 0, flag: "--force"},
		{name: "No Dependencies", otype: 0, flag: "--no-dependencies"},
		{name: "No Incremental", otype: 0, flag: "--no-incremental"},
		{name: "No Restore", otype: 0, flag: "--no-restore"},
		{name: "No Logo", otype: 0, flag: "--nologo"},
		{name: "No Self Contained", otype: 0, flag: "--no-self-contained"},
		{name: "Output", otype: 1, flag: "-o"},
		{name: "OS", otype: 1, flag: "--os"},
		{name: "Runtime", otype: 1, flag: "-r"},
		{name: "Source", otype: 1, flag: "--source"},
		{name: "Verbosity", otype: 1, flag: "-v"},
		{name: "Version Suffix", otype: 1, flag: "--version-suffix"},
		{name: "Confirm", otype: 2},
		{name: "Exit", otype: 2},
	}

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 7
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(options))

	return complexBuildModel{
		options:   options,
		paginator: p,
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
				} else if m.paginator.Page != 0 {
					m.cursor = m.paginator.PerPage
					m.paginator.Page--
				}

			case "down", "j":
				if m.cursor < m.paginator.ItemsOnPage(len(m.options))-1 {
					m.cursor++
				} else if m.paginator.Page != m.paginator.TotalPages-1 {
					m.cursor = 0
					m.paginator.Page++
				}

			case "left", "right", "h", "l":
				m.cursor = 0

			case "enter":
				switch m.options[m.cursor+m.paginator.Page*m.paginator.PerPage].otype {
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

		var cmd tea.Cmd
		m.paginator, cmd = m.paginator.Update(msg)
		return m, cmd

	case 1:
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k", "down", "j":
				m.optionBool = !m.optionBool
			case "enter":
				m.state = 0
				if m.optionBool {
					m.options[m.cursor+m.paginator.Page*m.paginator.PerPage].value = "Yes"
				} else {
					m.options[m.cursor+m.paginator.Page*m.paginator.PerPage].value = "No"
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
				m.options[m.cursor+m.paginator.Page*m.paginator.PerPage].value = m.input.Value()
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
		start, end := m.paginator.GetSliceBounds(len(m.options))
		for i, choice := range m.options[start:end] {
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

		s += "\n" + m.paginator.View()

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
