package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type InputPattern struct {
	DirName    string   `json:"dir_name"`
	Categories []string `json:"categories"`
}

type InputFile struct {
	Patterns []InputPattern `json:"patterns"`
}

type ExistingFile struct {
	DirName             string   `json:"dir_name"`
	FriendlyName        string   `json:"friendly_name"`
	ShortDescription    string   `json:"short_description"`
	Description         string   `json:"description"`
	Categories          []string `json:"categories"`
	Tags                []string `json:"tags"`
	RelatedPatterns     []string `json:"related_patterns"`
	CharacterCount      int      `json:"character_count"`
	EstimatedTokenCount int      `json:"estimated_token_count"`
	UsageExample        string   `json:"usage_example"`
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Get environment variables
	inputJSONPath := os.Getenv("JSON_UPDATES_PATH")
	metadataDir := os.Getenv("METADATA_DIR")

	// Read the input JSON file
	inputData, err := ioutil.ReadFile(inputJSONPath)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	var inputFile InputFile
	err = json.Unmarshal(inputData, &inputFile)
	if err != nil {
		fmt.Println("Error parsing input JSON:", err)
		return
	}

	// Process each pattern in the input file
	for _, pattern := range inputFile.Patterns {
		filename := filepath.Join(metadataDir, pattern.DirName+".json")
		
		// Read the existing JSON file
		existingData, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err)
			continue
		}

		var existingFile ExistingFile
		err = json.Unmarshal(existingData, &existingFile)
		if err != nil {
			fmt.Printf("Error parsing existing JSON in %s: %v\n", filename, err)
			continue
		}

		// Update the categories
		existingFile.Categories = pattern.Categories

		// Write the updated JSON back to the file
		updatedData, err := json.MarshalIndent(existingFile, "", "  ")
		if err != nil {
			fmt.Printf("Error creating updated JSON for %s: %v\n", filename, err)
			continue
		}

		err = ioutil.WriteFile(filename, updatedData, 0644)
		if err != nil {
			fmt.Printf("Error writing updated JSON to %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("Successfully updated %s\n", filename)
	}
}
