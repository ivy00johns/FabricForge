package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type Pattern struct {
	DirName      string   `json:"dir_name"`
	FriendlyName string   `json:"friendly_name"`
	ShortDesc    string   `json:"short_description"`
	Categories   []string `json:"categories"`
	Tags         []string `json:"tags"`
}

type PatternList struct {
	Patterns []Pattern `json:"patterns"`
}

type FilterOption struct {
	Name string
	Desc string
}

func (f FilterOption) Title() string       { return f.Name }
func (f FilterOption) Description() string { return f.Desc }
func (f FilterOption) FilterValue() string { return f.Name }

type model struct {
	list           list.Model
	textInput      textinput.Model
	allPatterns    []list.Item
	filteredItems  []list.Item
	alphaSort      bool
	sortByDirName  bool
	config         Config
	state          string
	selectedCmd    string
	confirmItems   []list.Item
	filterOptions  []list.Item
	currentFilter  string
	allTags        []string
	allCategories  []string
	allDirectories []string
}

func (i Pattern) Title() string {
	return fmt.Sprintf("%s (📂: \"%s\")", i.FriendlyName, i.DirName)
}

func (i Pattern) Description() string { return i.ShortDesc }
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

	filterOptions := []list.Item{
		FilterOption{Name: "Global Search", Desc: "Search across all fields"},
		FilterOption{Name: "Tags", Desc: "Filter by tags"},
		FilterOption{Name: "Categories", Desc: "Filter by categories"},
		FilterOption{Name: "Directories", Desc: "Filter by directory names"},
	}

	allTags, allCategories, allDirectories := extractMetadata(patterns)

	return model{
		list:           l,
		textInput:      ti,
		allPatterns:    patterns,
		filteredItems:  patterns,
		alphaSort:      config.AlphaSort,
		sortByDirName:  config.SortByDirName,
		config:         config,
		state:          "selecting",
		confirmItems:   confirmItems,
		filterOptions:  filterOptions,
		currentFilter:  "Global Search",
		allTags:        allTags,
		allCategories:  allCategories,
		allDirectories: allDirectories,
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

func extractMetadata(patterns []list.Item) ([]string, []string, []string) {
	tagMap := make(map[string]bool)
	categoryMap := make(map[string]bool)
	directoryMap := make(map[string]bool)

	for _, item := range patterns {
		pattern := item.(Pattern)
		for _, tag := range pattern.Tags {
			tagMap[tag] = true
		}
		for _, category := range pattern.Categories {
			categoryMap[category] = true
		}
		directoryMap[pattern.DirName] = true
	}

	return mapToSortedSlice(tagMap), mapToSortedSlice(categoryMap), mapToSortedSlice(directoryMap)
}

func mapToSortedSlice(m map[string]bool) []string {
	slice := make([]string, 0, len(m))
	for k := range m {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	return slice
}
