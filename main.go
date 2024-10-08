package main

import (
	"fmt"
	"os"

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
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
