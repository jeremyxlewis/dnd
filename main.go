/*
Copyright © 2025
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"dnd-cli/internal/data"
	"dnd-cli/internal/tui"
)

func main() {
	// Command line flags
	compendiumFile := flag.String("compendium", "", "Compendium XML file to load (default: Complete_Compendium.xml)")
	dataPath := flag.String("path", "./Compendium", "Path to Compendium directory")
	flag.Parse()

	// Use Compendium directory
	err := data.LoadData(*dataPath)
	if err != nil {
		fmt.Printf("Hark! The ancient scrolls of knowledge are sealed! Failed to load D&D data: %v\n", err)
		os.Exit(1)
	}

	// Load specific compendium file if provided
	if *compendiumFile != "" {
		err = data.LoadCompendiumData(*dataPath, *compendiumFile)
		if err != nil {
			fmt.Printf("Failed to load compendium file '%s': %v\n", *compendiumFile, err)
			os.Exit(1)
		}
	}

	// Start the TUI directly
	tui.StartTUI()
}
