package main

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type singleCommandModel struct {
	spinner    spinner.Model
	successMsg string
	ch         chan cmdOut
}

type cmdOut struct {
	err  error
	buff []byte
}

func createSingleCommandModel(cmd *exec.Cmd, msg string, s spinner.Model) singleCommandModel {
	ch := make(chan cmdOut)

	// Start console command
	go func() {
		// Create output bytes buffer and redirect cmd out
		var outbuff bytes.Buffer
		cmd.Stdout = &outbuff
		cmd.Stderr = &outbuff

		// Run cmd
		err := cmd.Run()

		// Send output via channel
		ch <- cmdOut{err: err, buff: outbuff.Bytes()}
	}()

	// Create model with configured channel
	return singleCommandModel{
		spinner:    s,
		successMsg: msg,
		ch:         ch,
	}
}

func (m singleCommandModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m singleCommandModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case cmd := <-m.ch:
		if cmd.err != nil {
			return createExitModel(
				fmt.Sprint(redItemStyle.Render("\u274c An error occured! "), "Exit log:\n\n", string(cmd.buff)),
			), tea.Quit
		}
		return createExitModel(greenItemStyle.Render("\u2705", m.successMsg)), tea.Quit

	default:
		switch msg := msg.(type) {

		case tea.KeyMsg:

			switch msg.String() {

			case "ctrl+c", "q":
				return createExitModel("Exiting..."), tea.Quit
			}
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m singleCommandModel) View() string {
	s := titleStyle.Render("DOTHELP") + "\n\n"

	s += fmt.Sprintf("%s Running process...\n", m.spinner.View())

	s += helpStyle.Render("\nPress q to quit.")

	return s
}
