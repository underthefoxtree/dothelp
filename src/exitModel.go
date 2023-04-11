package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type exitModel struct {
	exitMessage string
}

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
