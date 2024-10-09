package main

import (
	"encoding/json"
	"io/ioutil"
	"sort"

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

func sortPatterns(patterns []list.Item, alphaSort, sortByDirName bool) {
	sort.Slice(patterns, func(i, j int) bool {
		pi, pj := patterns[i].(Pattern), patterns[j].(Pattern)
		if sortByDirName {
			return pi.DirName < pj.DirName
		}
		if alphaSort {
			return pi.FriendlyName < pj.FriendlyName
		}
		return false
	})
}
