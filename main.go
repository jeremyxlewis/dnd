/*
Copyright © 2025
*/
package main

import (
	"fmt"
	"os"

	"dnd-cli/internal/data"
	"dnd-cli/internal/tui"
)

func main() {
	// The data files are located in data/ (relative to project root)
	dataPath := "./data"

	// Load the D&D data
	err := data.LoadData(dataPath)
	if err != nil {
		fmt.Printf("Hark! The ancient scrolls of knowledge are sealed! Failed to load D&D data: %v\n", err)
		os.Exit(1)
	}

	// Start the TUI directly
	tui.StartTUI()
}
