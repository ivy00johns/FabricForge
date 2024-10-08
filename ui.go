package main

import (
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
			filtered := filterPatterns(m.allPatterns, m.textInput.Value())
			sortPatterns(filtered, m.alphaSort, m.sortByDirName)
			m.list.SetItems(filtered)
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
	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(m.config.Title),
		"\n",
		"Search: "+m.textInput.View(),
		"\n",
		m.list.View(),
	))
}
