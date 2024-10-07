package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/charmbracelet/bubbles/list"
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
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

func (i Pattern) Title() string       { return i.FriendlyName }
func (i Pattern) Description() string { return i.ShortDescription }
func (i Pattern) FilterValue() string { return i.FriendlyName }

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	// Get metadata path from environment variable
	metadataPath := os.Getenv("MERGED_PATTERNS_METADATA_PATH")
	if metadataPath == "" {
		fmt.Println("MERGED_PATTERNS_METADATA_PATH not set in .env file")
		os.Exit(1)
	}

	// Load patterns from JSON file
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

	// Convert patterns to list items
	items := make([]list.Item, len(patternList.Patterns))
	for i, pattern := range patternList.Patterns {
		items[i] = pattern
	}

	// Create the list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Fabric Patterns"

	// Create the model
	m := model{list: l}

	// Run the program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
