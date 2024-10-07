package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

type Pattern struct {
	DirName          string   `json:"dir_name"`
	FriendlyName     string   `json:"friendly_name"`
	ShortDescription string   `json:"short_description"`
	Categories       []string `json:"categories"`
	Tags             []string `json:"tags"`
}

type PatternList struct {
	Patterns []Pattern `json:"patterns"`
}

type model struct {
	list        list.Model
	textInput   textinput.Model
	allPatterns []list.Item
}

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
			m.list.SetItems(filterPatterns(m.allPatterns, m.textInput.Value()))
		}

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-3)
		m.textInput.Width = msg.Width - h
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.textInput.View(),
		m.list.View(),
	)
}

func (i Pattern) Title() string       { return i.FriendlyName }
func (i Pattern) Description() string { return i.ShortDescription }
func (i Pattern) FilterValue() string { return i.FriendlyName }

func filterPatterns(patterns []list.Item, filter string) []list.Item {
	if filter == "" {
		return patterns
	}

	var filtered []list.Item
	for _, item := range patterns {
		pattern := item.(Pattern)
		if strings.Contains(strings.ToLower(pattern.FriendlyName), strings.ToLower(filter)) ||
			strings.Contains(strings.ToLower(pattern.ShortDescription), strings.ToLower(filter)) ||
			containsInSlice(pattern.Categories, filter) ||
			containsInSlice(pattern.Tags, filter) {
			filtered = append(filtered, pattern)
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

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	metadataPath := os.Getenv("MERGED_PATTERNS_METADATA_PATH")
	if metadataPath == "" {
		fmt.Println("MERGED_PATTERNS_METADATA_PATH not set in .env file")
		os.Exit(1)
	}

	file, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	var patternList PatternList
	err = json.Unmarshal(file, &patternList)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		os.Exit(1)
	}

	items := make([]list.Item, len(patternList.Patterns))
	for i, pattern := range patternList.Patterns {
		items[i] = pattern
	}

	ti := textinput.New()
	ti.Placeholder = "Type to filter patterns"
	ti.Focus()

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Fabric Patterns"

	m := model{
		list:        l,
		textInput:   ti,
		allPatterns: items,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
