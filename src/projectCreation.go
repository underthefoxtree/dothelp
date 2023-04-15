package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type projectCreationModel struct {
	templateName string
	templateId   string
	nameInput    textinput.Model
	outInput     textinput.Model
	cursor       int
}

func createProjectCreationModel(name string, id string) projectCreationModel {
	tiName := textinput.New()
	tiOut := textinput.New()

	name = strings.ReplaceAll(name, "\u00a0", " ")
	re := regexp.MustCompile("[^\x20-\x7E]+")
	name = re.ReplaceAllString(name, "")

	tiName.Placeholder = name
	tiOut.Placeholder = "./"
	tiName.Prompt = ""
	tiOut.Prompt = ""

	tiName.Focus()

	return projectCreationModel{
		templateName: name,
		templateId:   id,
		nameInput:    tiName,
		outInput:     tiOut,
	}
}

func (m projectCreationModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m projectCreationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c":
			return createExitModel("Exiting..."), tea.Quit

		case "up", "shift+tab":
			switch m.cursor {
			case 0:
				m.cursor = 3
				m.nameInput.Blur()
			case 1:
				m.cursor--
				m.nameInput.Focus()
				m.outInput.Blur()
			case 2:
				m.cursor--
				m.outInput.Focus()
			case 3:
				m.cursor--
			}

		case "down", "tab":
			switch m.cursor {
			case 0:
				m.cursor++
				m.nameInput.Blur()
				m.outInput.Focus()
			case 1:
				m.cursor++
				m.outInput.Blur()
			case 2:
				m.cursor++
			case 3:
				m.cursor = 0
				m.nameInput.Focus()
			}

		case "enter":
			switch m.cursor {
			case 0:
				m.cursor++
				m.nameInput.Blur()
				m.outInput.Focus()
			case 1:
				m.cursor++
				m.outInput.Blur()
			case 2:
				name, dir := m.templateName, "./"

				if m.outInput.Value() != "" {
					dir = m.outInput.Value()
				}
				if m.nameInput.Value() != "" {
					name = m.nameInput.Value()
				}

				return createExitModel(
					fmt.Sprintf("Creating %s in %s", name, dir),
				), tea.Quit
			case 3:
				return createExitModel("Exiting..."), tea.Quit
			}
		}
	}

	cmds := make([]tea.Cmd, 2)

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.outInput, cmds[1] = m.outInput.Update(msg)

	m.nameInput.Placeholder = m.templateName
	m.outInput.Placeholder = fmt.Sprintf("./%s", m.nameInput.Value())

	return m, tea.Batch(cmds...)
}

func (m projectCreationModel) View() string {
	s := titleStyle.Render("DOTHELP") + "\n\n"

	s += fmt.Sprintf("Creating %s.\n\n", m.templateName)

	if m.cursor == 0 {
		s += selectedItemStyle.Render("Name of project") + "\n"
		s += "> " + m.nameInput.View() + "\n"
	} else {
		s += itemStyle.Render("Name of project") + "\n"
		s += "  " + m.nameInput.View() + "\n"
	}

	if m.cursor == 1 {
		s += selectedItemStyle.Render("Output Directory") + "\n"
		s += "> " + m.outInput.View() + "\n"
	} else {
		s += itemStyle.Render("Output Directory") + "\n"
		s += "  " + m.outInput.View() + "\n"
	}

	s += "\n"

	if m.cursor == 2 {
		s += fmt.Sprintf("> %s", greenItemStyle.Render("Confirm"))
	} else {
		s += fmt.Sprintf("  %s", itemStyle.Render("Confirm"))
	}

	s += " | "

	if m.cursor == 3 {
		s += fmt.Sprintf("> %s", redItemStyle.Render("Exit"))
	} else {
		s += fmt.Sprintf("  %s", itemStyle.Render("Exit"))
	}

	s += helpStyle.Render("\n\nPress ctrl+c to quit.")

	return s
}
