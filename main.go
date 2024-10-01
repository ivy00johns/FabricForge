package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type Pattern struct {
	Name        string
	Description string
	Execute     func() error
}

var patterns = []Pattern{
	{
		Name:        "Pattern 1",
		Description: "Description of Pattern 1",
		Execute: func() error {
			fmt.Println("Executing Pattern 1")
			return nil
		},
	},
	{
		Name:        "Pattern 2",
		Description: "Description of Pattern 2",
		Execute: func() error {
			fmt.Println("Executing Pattern 2")
			return nil
		},
	},
	// Add more patterns here
}

func main() {
	for {
		pattern, err := selectPattern()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = pattern.Execute()
		if err != nil {
			fmt.Printf("Execution failed %v\n", err)
			return
		}

		continuePrompt := promptui.Select{
			Label: "Do you want to execute another pattern?",
			Items: []string{"Yes", "No"},
		}

		_, result, err := continuePrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if result == "No" {
			fmt.Println("Exiting...")
			break
		}
	}
}

func selectPattern() (Pattern, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Name | cyan }} ({{ .Description | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Description | red }})",
		Selected: "\U0001F336 {{ .Name | red | cyan }}",
		Details: `
--------- Pattern ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Description:" | faint }}	{{ .Description }}`,
	}

	prompt := promptui.Select{
		Label:     "Select Pattern",
		Items:     patterns,
		Templates: templates,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return Pattern{}, err
	}

	return patterns[i], nil
}
