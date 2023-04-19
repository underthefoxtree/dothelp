package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func createSingleCommandModel(cmd *exec.Cmd, msg string) (singleCommandModel, tea.Cmd) {
	s := spinner.New()

	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

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
	}, s.Tick
}

func clearScreen() string {
	return "\r\033[2K\x1b[1A\r\033[2K\x1b[1A\r\033[2K\x1b[1A\r\033[2K\x1b[1A\r"
}

func (m singleCommandModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m singleCommandModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case cmd := <-m.ch:
		if cmd.err != nil {
			fmt.Print(clearScreen(), redItemStyle.Render("\u274c An error occured! "), "Exit log:\n\n", string(cmd.buff))
			os.Exit(1)
		}
		return createExitModel(greenItemStyle.Render("\u2705", m.successMsg))

	default:
		switch msg := msg.(type) {

		case tea.KeyMsg:

			switch msg.String() {

			case "ctrl+c", "q":
				return createExitModel("Exiting...")
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
