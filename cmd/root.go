/*
Copyright © 2025
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dnd-cli/internal/data"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dnd",
	Short: "A CLI companion for Dungeons & Dragons",
	Long: `dnd is a command-line companion for Dungeons & Dragons players and Dungeon Masters.
It provides tools for rolling dice, looking up rules, spells, monsters, and more.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Determine the path to the data directory.
	// This assumes the 'data' repository is cloned directly into the project root
	// and contains a 'data' subdirectory with the JSON files.
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Hark! A shadow falls upon our path: Failed to get current working directory: %v\n", err)
		os.Exit(1)
	}

	// The data files are located in dnd/data/data/ (relative to project root)
	dataPath := filepath.Join(currentDir, "data", "data")

	// Load the D&D data
	err = data.LoadData(dataPath)
	if err != nil {
		fmt.Printf("Hark! The ancient scrolls of knowledge are sealed! Failed to load D&D data: %v\n", err)
		os.Exit(1)
	}

	// Add commands
	// RootCmd.AddCommand(charCmd) // Removed for now
}
