package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

type PatternMetadata struct {
	DirName      string   `json:"dir_name"`
	FriendlyName string   `json:"friendly_name"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
}

type PatternsMetadata struct {
	Patterns []PatternMetadata `json:"patterns"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	patternsDir := os.Getenv("FABRIC_PATTERNS_DIRECTORY_PATH")
	if patternsDir == "" {
		fmt.Println("FABRIC_PATTERNS_DIRECTORY_PATH not set in .env file")
		return
	}

	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		fmt.Println("OUTPUT_DIR not set in .env file")
		return
	}

	metadataPath := os.Getenv("PATTERNS_METADATA_PATH")
	if metadataPath == "" {
		fmt.Println("PATTERNS_METADATA_PATH not set in .env file")
		return
	}

	patterns, err := loadPatternMetadata(metadataPath)
	if err != nil {
		fmt.Printf("Error loading pattern metadata: %v\n", err)
		return
	}

	for {
		pattern, err := selectPattern(patterns)
		if err != nil {
			fmt.Printf("Pattern selection failed: %v\n", err)
			return
		}

		inputSource, err := selectInputSource()
		if err != nil {
			fmt.Printf("Input source selection failed: %v\n", err)
			return
		}

		command, err := buildFabricCommand(pattern.DirName, inputSource, outputDir)
		if err != nil {
			fmt.Printf("Command building failed: %v\n", err)
			return
		}

		fmt.Printf("\nCommand to be executed:\n%s\n\n", command)

		if !confirmExecution() {
			fmt.Println("Execution cancelled.")
			if !promptContinue() {
				break
			}
			continue
		}

		err = executeFabricCommand(command)
		if err != nil {
			fmt.Printf("Execution failed: %v\n", err)
			return
		}

		if !promptContinue() {
			break
		}
	}
}

func loadPatternMetadata(path string) ([]PatternMetadata, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata PatternsMetadata
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil, err
	}

	sort.Slice(metadata.Patterns, func(i, j int) bool {
		return metadata.Patterns[i].Category < metadata.Patterns[j].Category
	})

	return metadata.Patterns, nil
}

func selectPattern(patterns []PatternMetadata) (PatternMetadata, error) {
	funcMap := template.FuncMap{
		"join":  strings.Join,
		"cyan":  color.CyanString,
		"green": color.GreenString,
		"red":   color.RedString,
		"faint": color.New(color.Faint).SprintFunc(),
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ .Category | cyan }}",
		Active:   "\U0001F4C1 {{ .FriendlyName | cyan }} ({{ .Description | green }})",
		Inactive: "  {{ .FriendlyName | cyan }} ({{ .Description | green }})",
		Selected: "\U0001F4C1 {{ .FriendlyName | red | cyan }}",
		Details: `
--------- Pattern ----------
{{ "Name:" | faint }}	{{ .FriendlyName }}
{{ "Description:" | faint }}	{{ .Description }}
{{ "Category:" | faint }}	{{ .Category }}
{{ "Tags:" | faint }}	{{ join .Tags ", " }}`,
		FuncMap: funcMap,
	}

	prompt := promptui.Select{
		Label:     "Select Pattern",
		Items:     patterns,
		Templates: templates,
		Size:      10,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return PatternMetadata{}, err
	}

	return patterns[i], nil
}

func selectInputSource() (string, error) {
	prompt := promptui.Select{
		Label: "Select input source",
		Items: []string{"Clipboard", "Manual Input"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func buildFabricCommand(pattern, inputSource, outputDir string) (string, error) {
	streamResults := os.Getenv("STREAM_RESULTS") == "true"
	streamFlag := ""
	if streamResults {
		streamFlag = "--stream"
	}

	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_output.txt", pattern))

	var command string
	if inputSource == "Clipboard" {
		command = fmt.Sprintf("pbpaste | fabric %s --pattern %s > %s", streamFlag, pattern, outputFile)
	} else {
		prompt := promptui.Prompt{
			Label: "Enter your text",
		}
		input, err := prompt.Run()
		if err != nil {
			return "", err
		}
		command = fmt.Sprintf("echo '%s' | fabric %s --pattern %s > %s", input, streamFlag, pattern, outputFile)
	}

	return command, nil
}

func confirmExecution() bool {
	prompt := promptui.Select{
		Label: "Do you want to execute this command?",
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return false
	}

	return result == "Yes"
}

func executeFabricCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("Command executed successfully.\n")
	return nil
}

func promptContinue() bool {
	prompt := promptui.Select{
		Label: "Do you want to execute another pattern?",
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return false
	}

	return result == "Yes"
}
