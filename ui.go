package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	appStyle = lipgloss.NewStyle().
			Padding(1, 2, 1, 2)
)

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.state {
			case "selecting":
				if m.list.SelectedItem() != nil {
					selectedPattern := m.list.SelectedItem().(Pattern)
					m.selectedCmd = m.buildFabricCommand(selectedPattern)
					m.state = "confirming"
					m.list.SetItems(m.confirmItems)
				}
			case "confirming":
				if m.list.SelectedItem() != nil {
					choice := m.list.SelectedItem().(confirmItem)
					if choice.title == "Yes" {
						m.state = "executing"
						return m, tea.Quit
					} else {
						m.state = "selecting"
						m.list.SetItems(m.allPatterns)
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-6)
		m.textInput.Width = msg.Width - h - 4
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	var content string
	switch m.state {
	case "selecting":
		content = lipgloss.JoinVertical(lipgloss.Left,
			"Select a pattern (↑/↓ to navigate, enter to select):",
			m.list.View(),
		)
	case "confirming":
		content = lipgloss.JoinVertical(lipgloss.Left,
			fmt.Sprintf("Command to execute: %s", m.selectedCmd),
			"Do you want to execute this command?",
			m.list.View(),
		)
	}

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(m.config.Title),
		"\n",
		content,
	))
}
