package main

import (
	"os"
	"strconv"
)

type Config struct {
	Width         int
	Height        int
	Title         string
	Placeholder   string
	MetadataPath  string
	AlphaSort     bool
	SortByDirName bool
	OutputDir     string
	StreamResults bool
	OutputResults bool
}

func loadConfig() Config {
	width, _ := strconv.Atoi(os.Getenv("CLI_WIDTH"))
	height, _ := strconv.Atoi(os.Getenv("CLI_HEIGHT"))
	alphaSort, _ := strconv.ParseBool(os.Getenv("ALPHA_SORT"))
	sortByDirName, _ := strconv.ParseBool(os.Getenv("SORT_BY_DIR_NAME"))
	streamResults, _ := strconv.ParseBool(os.Getenv("STREAM_RESULTS"))
	outputResults, _ := strconv.ParseBool(os.Getenv("OUTPUT_RESULTS"))

	return Config{
		Width:         width,
		Height:        height,
		Title:         os.Getenv("CLI_TITLE"),
		Placeholder:   os.Getenv("CLI_PLACEHOLDER"),
		MetadataPath:  os.Getenv("MERGED_PATTERNS_METADATA_PATH"),
		AlphaSort:     alphaSort,
		SortByDirName: sortByDirName,
		OutputDir:     os.Getenv("OUTPUT_DIR"),
		StreamResults: streamResults,
		OutputResults: outputResults,
	}
}
