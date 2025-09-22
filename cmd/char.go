package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"dnd-cli/internal/character"
	"dnd-cli/internal/data"

	"github.com/spf13/cobra"
)

// charCmd represents the char command
var charCmd = &cobra.Command{
	Use:   "char",
	Short: "Manage D&D characters (create, view, levelup)",
	Long: `The char command provides subcommands to manage your D&D characters.

Use 'dnd char create <name>' to make a new character.
Use 'dnd char view <name>' to see a character's sheet.
Use 'dnd char levelup <name>' to advance a character's level.`,
}

func init() {
	RootCmd.AddCommand(charCmd)

	// Add 'create' subcommand
	var createCharCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new D&D character",
		Long:  `Guides you through the process of creating a new D&D character.`, 
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]

			// Check if character already exists
			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}
			if _, err := os.Stat(charFilePath); err == nil {
				fmt.Printf("Hark! A hero named '%s' already exists in the annals! Choose another name, adventurer.\n", charName)
				return
			}

			reader := bufio.NewReader(os.Stdin)

			// Species selection
			fmt.Printf("Choose a species (e.g., Human, Elf, Dwarf). Available: %s\n", getSpeciesNames())
			fmt.Print("Enter Species: ")
			speciesInput, _ := reader.ReadString('\n')
			speciesInput = strings.TrimSpace(speciesInput)
			if _, err := data.GetSpeciesByName(speciesInput); err != nil {
				fmt.Printf("Hark! That species is unknown to these lands. %v\n", err)
				return
			}

			// Background selection
			fmt.Printf("Choose a background (e.g., Acolyte, Soldier, Criminal). Available: %s\n", getBackgroundNames())
			fmt.Print("Enter Background: ")
			backgroundInput, _ := reader.ReadString('\n')
			backgroundInput = strings.TrimSpace(backgroundInput)
			if _, err := data.GetBackgroundByName(backgroundInput); err != nil {
				fmt.Printf("Hark! That background is not etched in our lore. %v\n", err)
				return
			}

			// Simplified ability score generation (all 10 for now)
			newChar := character.NewCharacter(charName, speciesInput, "Adventurer", backgroundInput) // Class is placeholder

			// Save character
			err = character.SaveCharacter(newChar, charFilePath)
			if err != nil {
				fmt.Printf("Hark! The scribe's quill faltered! Failed to save character: %v\n", err)
				return
			}

			fmt.Printf("\nVerily! Character '%s' created and saved to '%s'.\n", charName, charFilePath)
		},
	}
	charCmd.AddCommand(createCharCmd)

	// Add 'view' subcommand
	var viewCharCmd = &cobra.Command{
		Use:   "view [name]",
		Short: "View a D&D character's sheet",
		Long:  `Displays the details of a saved D&D character.`, 
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]

			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}

			char, err := character.LoadCharacter(charFilePath)
			if err != nil {
				fmt.Printf("Hark! The hero '%s' is not found in the archives! %v\n", charName, err)
				return
			}

			fmt.Printf("\n--- Character Sheet: %s ---\n", char.Name)
			fmt.Printf("Species: %s\n", char.Species)
			fmt.Printf("Class: %s\n", char.Class)
			fmt.Printf("Level: %d\n", char.Level)
			fmt.Printf("Background: %s\n", char.Background)
			fmt.Println("\n--- Ability Scores ---")
			fmt.Printf("STR: %d\n", char.Strength)
			fmt.Printf("DEX: %d\n", char.Dexterity)
			fmt.Printf("CON: %d\n", char.Constitution)
			fmt.Printf("INT: %d\n", char.Intelligence)
			fmt.Printf("WIS: %d\n", char.Wisdom)
			fmt.Printf("CHA: %d\n", char.Charisma)
			fmt.Println("\n--- Derived Stats ---")
			fmt.Printf("HP: %d\n", char.HitPoints)
			fmt.Printf("AC: %d\n", char.ArmorClass)
			fmt.Printf("Proficiency Bonus: +%d\n", char.ProficiencyBonus)
			fmt.Println("---------------------------\n")
		},
	}
	charCmd.AddCommand(viewCharCmd)

	// Add 'levelup' subcommand
	var levelUpCharCmd = &cobra.Command{
		Use:   "levelup [name]",
		Short: "Level up a D&D character",
		Long:  `Increases the level of a saved D&D character and updates basic stats.`, 
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]

			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}

			char, err := character.LoadCharacter(charFilePath)
			if err != nil {
				fmt.Printf("Hark! The hero '%s' is not found in the archives! %v\n", charName, err)
				return
			}

			fmt.Printf("\nVerily! '%s' is now level %d.\n", char.Name, char.Level+1)
			char.LevelUp()

			err = character.SaveCharacter(char, charFilePath)
			if err != nil {
				fmt.Printf("Hark! The scribe's quill faltered! Failed to save character: %v\n", err)
				return
			}

			fmt.Printf("Character '%s' leveled up to %d! HP: %d, Proficiency Bonus: +%d.\n", char.Name, char.Level, char.HitPoints, char.ProficiencyBonus)
		},
	}
	charCmd.AddCommand(levelUpCharCmd)
}

// Helper to get available species names
func getSpeciesNames() string {
	names := make([]string, len(data.AllSpecies))
	for i, s := range data.AllSpecies {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}

// Helper to get available background names
func getBackgroundNames() string {
	names := make([]string, len(data.AllBackgrounds))
	for i, b := range data.AllBackgrounds {
		names[i] = b.Name
	}
	return strings.Join(names, ", ")
}
