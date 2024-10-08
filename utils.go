package main

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

func loadPatterns(metadataPath string) ([]list.Item, error) {
	file, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var patternList PatternList
	err = json.Unmarshal(file, &patternList)
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, len(patternList.Patterns))
	for i, pattern := range patternList.Patterns {
		items[i] = pattern
	}

	return items, nil
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

func sortPatterns(patterns []list.Item, alphaSort, sortByDirName bool) {
	sort.Slice(patterns, func(i, j int) bool {
		pi, pj := patterns[i].(Pattern), patterns[j].(Pattern)
		if sortByDirName {
			return pi.DirName < pj.DirName
		}
		if alphaSort {
			return pi.FriendlyName < pj.FriendlyName
		}
		return false // No sorting
	})
}
