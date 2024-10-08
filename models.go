package main

import (
	"fmt"

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
}

func (i Pattern) Title() string {
	return fmt.Sprintf("%s (📂: \"%s\")", i.FriendlyName, i.DirName)
}

func (i Pattern) Description() string { return i.ShortDescription }
func (i Pattern) FilterValue() string { return i.FriendlyName }

func initialModel(patterns []list.Item, config Config) model {
	ti := textinput.New()
	ti.Placeholder = config.Placeholder
	ti.Focus()

	l := list.New(patterns, list.NewDefaultDelegate(), config.Width, config.Height)

	return model{
		list:          l,
		textInput:     ti,
		allPatterns:   patterns,
		alphaSort:     config.AlphaSort,
		sortByDirName: config.SortByDirName,
		config:        config,
	}
}
