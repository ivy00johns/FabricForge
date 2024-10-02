package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Metadata represents the structure of each JSON file's metadata
type Metadata struct {
	DirName            string   `json:"dir_name"`
	FriendlyName       string   `json:"friendly_name"`
	ShortDescription   string   `json:"short_description"`
	Description        string   `json:"description"`
	Categories         []string `json:"categories"`
	Tags               []string `json:"tags"`
	RelatedPatterns    []string `json:"related_patterns"`
	CharacterCount     int      `json:"character_count"`
	EstimatedTokenCount int     `json:"estimated_token_count"`
	UsageExample       string   `json:"usage_example"`
}

// CombinedMetadata holds the collection of all metadata
type CombinedMetadata struct {
	Patterns []Metadata `json:"patterns"`
}

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve paths from the .env file
	metadataDir := os.Getenv("METADATA_DIR")
	outputDir := os.Getenv("OUTPUT_DIR")
	mergedMetadataFilePath := os.Getenv("MERGED_PATTERNS_METADATA_PATH")

	// Prepare to collect all metadata
	var combinedMetadata CombinedMetadata

	// Walk through the metadata directory and process .json files
	err = filepath.Walk(metadataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a .json file
		if filepath.Ext(path) == ".json" {
			fmt.Printf("Processing file: %s\n", path)  // Log the file name being processed

			// Read the file content
			fileContent, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Unmarshal the file content into a Metadata object
			var metadata Metadata
			if err := json.Unmarshal(fileContent, &metadata); err != nil {
				fmt.Printf("Error unmarshaling file %s: %v\n", path, err)
				return err
			}

			// Append the metadata to the combined collection
			combinedMetadata.Patterns = append(combinedMetadata.Patterns, metadata)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Check if any metadata was collected
	if len(combinedMetadata.Patterns) == 0 {
		fmt.Println("No metadata was merged. Please check the directory or file structure.")
		return
	}

	// Create the output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating output directory: %v", err)
		}
	}

	// Marshal the combined metadata to JSON with indentation for readability
	outputFile, err := json.MarshalIndent(combinedMetadata, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write the merged metadata to the specified output file
	if err := ioutil.WriteFile(mergedMetadataFilePath, outputFile, 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("All metadata merged successfully into %s\n", mergedMetadataFilePath)
}
