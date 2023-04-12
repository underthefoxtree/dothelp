package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type templateSelectModel struct {
	options   []template
	cursor    int
	paginator paginator.Model
}

type template struct {
	name string
	id   string
}

func createTemplateSelectModel() templateSelectModel {
	// Run dotnet new --list and capture output
	out, err := exec.Command("dotnet", "new", "--list").Output()
	if err != nil {
		panic(err)
	}

	// Vars to create template array
	var (
		templates            []template
		nameLength, idLength int
	)

	// Use scanner to read output line by line
	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	for scanner.Scan() {
		line := scanner.Text()

		// If the line starts with a -, use the spaces between to
		// get the maximum length of the names and ids
		if strings.HasPrefix(line, "-") {
			nameLength = strings.IndexByte(line, ' ')
			idLength = nameLength + 2 + strings.IndexByte(line[nameLength+2:], ' ')
		} else if nameLength != 0 && len(line) > nameLength*2 {
			// Add template using information
			templates = append(templates, template{
				name: strings.TrimSpace(line[:nameLength]),
				id:   strings.TrimSpace(line[nameLength+2 : idLength]),
			})
		}
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 7
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(templates))

	return templateSelectModel{
		options:   templates,
		cursor:    0,
		paginator: p,
	}
}

func (m templateSelectModel) Init() tea.Cmd {
	return nil
}

func (m templateSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return createExitModel("Exiting..."), tea.Quit

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

		case "enter", " ":
			return createExitModel(
					fmt.Sprintf(
						"You selected: %s",
						selectedItemStyle.Render(m.options[m.cursor+m.paginator.Page*m.paginator.PerPage].name))),
				tea.Quit
		}
	}

	var cmd tea.Cmd
	m.paginator, cmd = m.paginator.Update(msg)

	return m, cmd
}

func (m templateSelectModel) View() string {
	s := titleStyle.Render("DOTHELP") + "\n\n"

	start, end := m.paginator.GetSliceBounds(len(m.options))

	for i, choice := range m.options[start:end] {
		cursor := "  "
		style := itemStyle
		if i == m.cursor {
			cursor = "> "
			style = selectedItemStyle
		}

		s += fmt.Sprintf("%s%s (%s)\n", cursor, style.Render(choice.name), helpStyle.Render(choice.id))
	}

	s += "\n" + m.paginator.View()
	s += helpStyle.Render("\nPress q to quit.")

	return s
}
