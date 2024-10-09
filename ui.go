package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 0)

	appStyle = lipgloss.NewStyle().
			Padding(0, 0, 0, 0)
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
		case "/":
			if m.state == "selecting" {
				m.state = "filter_menu"
				m.list.SetItems(m.filterOptions)
			}
		case "esc":
			if m.state == "filtering" || m.state == "filter_menu" {
				m.state = "selecting"
				m.list.SetItems(m.filteredItems)
			}
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
						m.list.SetItems(m.filteredItems)
					}
				}
			case "filter_menu":
				if m.list.SelectedItem() != nil {
					filterOption := m.list.SelectedItem().(FilterOption)
					m.currentFilter = filterOption.Name
					m.state = "filtering"
					m.textInput.SetValue("")
					m.textInput.Focus()
					if m.currentFilter != "Global Search" {
						var filterItems []list.Item
						switch m.currentFilter {
						case "Tags":
							filterItems = stringSliceToListItems(m.allTags)
						case "Categories":
							filterItems = stringSliceToListItems(m.allCategories)
						case "Directories":
							filterItems = stringSliceToListItems(m.allDirectories)
						}
						m.list.SetItems(filterItems)
					}
				}
			case "filtering":
				if m.currentFilter != "Global Search" && m.list.SelectedItem() != nil {
					selectedFilter := m.list.SelectedItem().(FilterOption).Name
					m.filteredItems = filterPatternsByMetadata(m.allPatterns, m.currentFilter, selectedFilter)
					m.state = "selecting"
					m.list.SetItems(m.filteredItems)
				} else {
					m.filteredItems = filterPatterns(m.allPatterns, m.textInput.Value())
					m.state = "selecting"
					m.list.SetItems(m.filteredItems)
				}
			}
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-6)
		m.textInput.Width = msg.Width - h - 4
	}

	if m.state == "filtering" && m.currentFilter == "Global Search" {
		m.filteredItems = filterPatterns(m.allPatterns, m.textInput.Value())
		m.list.SetItems(m.filteredItems)
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
			"Select a pattern (↑/↓ to navigate, enter to select, / to filter):",
			m.list.View(),
		)
	case "confirming":
		content = lipgloss.JoinVertical(lipgloss.Left,
			fmt.Sprintf("Command to execute: %s\n", m.selectedCmd),
			"Do you want to execute this command?",
			m.list.View(),
		)
	case "filter_menu":
		content = lipgloss.JoinVertical(lipgloss.Left,
			"Select filter type:",
			m.list.View(),
		)
	case "filtering":
		if m.currentFilter == "Global Search" {
			content = lipgloss.JoinVertical(lipgloss.Left,
				fmt.Sprintf("Filter %s (type to filter, esc to cancel):", m.currentFilter),
				m.textInput.View(),
				m.list.View(),
			)
		} else {
			content = lipgloss.JoinVertical(lipgloss.Left,
				fmt.Sprintf("Select %s (↑/↓ to navigate, enter to select, esc to cancel):", m.currentFilter),
				m.list.View(),
			)
		}
	}

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		"\n",
		content,
	))
}

func filterPatterns(patterns []list.Item, filter string) []list.Item {
	if filter == "" {
		return patterns
	}

	var filtered []list.Item
	for _, item := range patterns {
		pattern := item.(Pattern)
		if strings.Contains(strings.ToLower(pattern.FriendlyName), strings.ToLower(filter)) ||
			strings.Contains(strings.ToLower(pattern.DirName), strings.ToLower(filter)) ||
			strings.Contains(strings.ToLower(pattern.ShortDesc), strings.ToLower(filter)) ||
			containsInSlice(pattern.Categories, filter) ||
			containsInSlice(pattern.Tags, filter) {
			filtered = append(filtered, pattern)
		}
	}
	return filtered
}

func filterPatternsByMetadata(patterns []list.Item, filterType, filterValue string) []list.Item {
	var filtered []list.Item
	for _, item := range patterns {
		pattern := item.(Pattern)
		switch filterType {
		case "Tags":
			if containsInSlice(pattern.Tags, filterValue) {
				filtered = append(filtered, pattern)
			}
		case "Categories":
			if containsInSlice(pattern.Categories, filterValue) {
				filtered = append(filtered, pattern)
			}
		case "Directories":
			if pattern.DirName == filterValue {
				filtered = append(filtered, pattern)
			}
		}
	}
	return filtered
}

func containsInSlice(slice []string, search string) bool {
	for _, item := range slice {
		if strings.Contains(strings.ToLower(item), strings.ToLower(search)) {
			return true
		}
	}
	return false
}

func stringSliceToListItems(slice []string) []list.Item {
	items := make([]list.Item, len(slice))
	for i, s := range slice {
		items[i] = FilterOption{Name: s, Desc: ""}
	}
	return items
}
