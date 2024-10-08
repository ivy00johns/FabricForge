package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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
	list          list.Model
	textInput     textinput.Model
	allPatterns   []list.Item
	alphaSort     bool
	sortByDirName bool
	config        Config
	state         string
	selectedCmd   string
	confirmItems  []list.Item
}

func (i Pattern) Title() string {
	return fmt.Sprintf("%s (📂: \"%s\")", i.FriendlyName, i.DirName)
}

func (i Pattern) Description() string { return i.ShortDescription }
func (i Pattern) FilterValue() string { return i.FriendlyName }

type confirmItem struct {
	title string
	desc  string
}

func (i confirmItem) Title() string       { return i.title }
func (i confirmItem) Description() string { return i.desc }
func (i confirmItem) FilterValue() string { return i.title }

func initialModel(patterns []list.Item, config Config) model {
	ti := textinput.New()
	ti.Placeholder = config.Placeholder
	ti.Focus()

	l := list.New(patterns, list.NewDefaultDelegate(), config.Width, config.Height)

	confirmItems := []list.Item{
		confirmItem{title: "Yes", desc: "Execute the command"},
		confirmItem{title: "No", desc: "Cancel and return to pattern selection"},
	}

	return model{
		list:          l,
		textInput:     ti,
		allPatterns:   patterns,
		alphaSort:     config.AlphaSort,
		sortByDirName: config.SortByDirName,
		config:        config,
		state:         "selecting",
		confirmItems:  confirmItems,
	}
}

func (m *model) buildFabricCommand(pattern Pattern) string {
	timestamp := time.Now().Format("2006-01-02T15:04:05-07:00")
	outputFile := fmt.Sprintf("%s_%s_output.md", pattern.DirName, timestamp)

	var command string
	if m.config.StreamResults && m.config.OutputResults {
		command = fmt.Sprintf("pbpaste | fabric --pattern %s | tee %s/%s", pattern.DirName, m.config.OutputDir, outputFile)
	} else if m.config.OutputResults {
		command = fmt.Sprintf("pbpaste | fabric --pattern %s > %s/%s", pattern.DirName, m.config.OutputDir, outputFile)
	} else {
		command = fmt.Sprintf("pbpaste | fabric --pattern %s", pattern.DirName)
	}

	return command
}
