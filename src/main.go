package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	config := loadConfig()
	patterns, err := loadPatterns(config.MetadataPath)
	if err != nil {
		fmt.Printf("Error loading patterns: %v\n", err)
		os.Exit(1)
	}

	m := initialModel(patterns, config)

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}

	// Execute the selected command only if confirmed
	finalM, ok := finalModel.(model)
	if ok && finalM.selectedCmd != "" && finalM.state == "executing" {
		fmt.Printf("Executing command: %s\n", finalM.selectedCmd)
		cmd := exec.Command("sh", "-c", finalM.selectedCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			os.Exit(1)
		}
	}
}