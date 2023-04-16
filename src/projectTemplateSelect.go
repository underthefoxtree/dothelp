package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type templateSelectModel struct {
	options   []template
	cursor    int
	paginator paginator.Model
	mode      mode
	filter    string
	input     textinput.Model
	filtered  []string
}

type template struct {
	name string
	id   string
}

type mode int

const (
	normal mode = iota
	search
)

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

	ti := textinput.New()
	ti.Placeholder = "Filter"
	ti.CharLimit = 60
	ti.Width = 30
	ti.Prompt = "  "

	return templateSelectModel{
		options:   templates,
		cursor:    0,
		paginator: p,
		mode:      normal,
		input:     ti,
	}
}

func (m templateSelectModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m templateSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case normal:
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

			case "/":
				m.input.Focus()
				m.mode = search

			case "enter", " ":
				t := m.options[m.cursor+m.paginator.Page*m.paginator.PerPage]

				return createProjectCreationModel(
					t.name,
					t.id,
				), nil
			}
		}

		var cmd tea.Cmd
		m.paginator, cmd = m.paginator.Update(msg)

		return m, cmd

	case search:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {

			case "ctrl+c", "q":
				return createExitModel("Exiting..."), tea.Quit

			case "up":
				if m.cursor > 0 {
					m.cursor--
				} else if m.paginator.Page != 0 {
					m.cursor = m.paginator.PerPage
					m.paginator.Page--
				}

			case "down":
				if m.cursor < m.paginator.ItemsOnPage(len(m.filtered))-1 {
					m.cursor++
				} else if m.paginator.Page != m.paginator.TotalPages-1 {
					m.cursor = 0
					m.paginator.Page++
				}

			case "left", "right", "h", "l":
				m.cursor = 0

			case "/", "esc":
				m.input.Blur()
				m.input.Reset()
				m.mode = normal

			case "enter":
				s := m.filtered[m.cursor+m.paginator.Page*m.paginator.PerPage]
				i := strings.IndexRune(s, '(')

				new := createProjectCreationModel(
					strings.TrimSpace(s[:i]),
					strings.TrimSpace(s[i+1:strings.LastIndex(s, ")")]),
				)

				return new, nil

			default:
				m.paginator.Page = 0
			}
		}

		var l []string

		for _, choice := range m.options {
			l = append(l, fmt.Sprintf("%s (%s)", choice.name, choice.id))
		}

		f := fuzzy.RankFindNormalizedFold(m.filter, l)
		sort.Sort(f)

		m.filtered = []string{}

		for _, i := range f {
			m.filtered = append(m.filtered, i.Target)
		}

		le := len(m.filtered) - 1

		if m.cursor > le {
			m.cursor = le
		}

		cmds := make([]tea.Cmd, 2)
		m.paginator.SetTotalPages(len(m.filtered))
		m.paginator, cmds[0] = m.paginator.Update(msg)
		m.input, cmds[1] = m.input.Update(msg)

		m.filter = m.input.Value()

		return m, tea.Batch(cmds...)

	default:
		return m, nil
	}
}

func (m templateSelectModel) RenderListNormal() string {
	var s string

	start, end := m.paginator.GetSliceBounds(len(m.options))

	for i, choice := range m.options[start:end] {
		cursor := "  "
		style := itemStyle
		if i == m.cursor {
			cursor = "> "
			style = selectedItemStyle
		}

		s += fmt.Sprintf("%s%s (%s)\n", cursor, style.Render(choice.name), style.Render(choice.id))
	}

	return s
}

func (m templateSelectModel) RenderListFiltered() string {
	var s string

	start, end := m.paginator.GetSliceBounds(len(m.filtered))

	for i, choice := range m.filtered[start:end] {
		cursor := "  "
		style := itemStyle
		if i == m.cursor {
			cursor = "> "
			style = selectedItemStyle
		}

		s += fmt.Sprintf("%s%s\n", cursor, style.Render(choice))
	}

	return s
}

func (m templateSelectModel) View() string {
	s := titleStyle.Render("DOTHELP") + "\n\n"

	if m.mode == search {
		s += m.input.View() + "\n"

		if len(m.filter) != 0 {
			s += m.RenderListFiltered()
		} else {
			s += m.RenderListNormal()
		}
	} else {
		s += m.RenderListNormal()
	}

	s += "\n" + m.paginator.View()
	s += helpStyle.Render("\nPress q to quit, / to filter.")

	return s
}
