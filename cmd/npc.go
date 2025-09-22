package cmd

import (
	"fmt"

	"dnd-cli/internal/data"

	"github.com/spf13/cobra"
)

// npcCmd represents the npc command
var npcCmd = &cobra.Command{
	Use:   "npc",
	Short: "Generates random NPCs with detailed backstories",
	Long: `The npc command generates a random non-player character (NPC) 
with a name, species, background, personality traits, ideals, bonds, flaws, and a backstory snippet to help DMs populate their world quickly.

Examples:
  dnd npc generate
  dnd npc`, // 'dnd npc' will default to generate
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is given, default to generate
		if len(args) == 0 || args[0] == "generate" {
			npc := data.GenerateNPC()
			fmt.Printf("\n--- Generated NPC ---\n")
			fmt.Printf("Name: %s\n", npc.Name)
			fmt.Printf("Species: %s\n", npc.Species)
			fmt.Printf("Background: %s\n", npc.Background)
			fmt.Printf("Personality Trait: %s\n", npc.PersonalityTrait)
			fmt.Printf("Ideal: %s\n", npc.Ideal)
			fmt.Printf("Bond: %s\n", npc.Bond)
			fmt.Printf("Flaw: %s\n", npc.Flaw)
			fmt.Printf("Backstory: %s\n", npc.Backstory)
			fmt.Println("---------------------\n")
		} else {
			fmt.Printf("Hark! Thy subcommand '%s' for npc doth bewilder my arcane senses. Pray tell, couldst thou rephrase thy plea.\n", args[0])
		}
	},
}

func init() {
	RootCmd.AddCommand(npcCmd)

	// Add a 'generate' subcommand explicitly for clarity, though 'dnd npc' will default to it.
	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generates a random NPC",
		Long:  `Generates a random non-player character (NPC) with a name, species, background, personality traits, ideals, bonds, flaws, and a backstory snippet.`,
		Run: func(cmd *cobra.Command, args []string) {
			npc := data.GenerateNPC()
			fmt.Printf("\n--- Generated NPC ---\n")
			fmt.Printf("Name: %s\n", npc.Name)
			fmt.Printf("Species: %s\n", npc.Species)
			fmt.Printf("Background: %s\n", npc.Background)
			fmt.Printf("Personality Trait: %s\n", npc.PersonalityTrait)
			fmt.Printf("Ideal: %s\n", npc.Ideal)
			fmt.Printf("Bond: %s\n", npc.Bond)
			fmt.Printf("Flaw: %s\n", npc.Flaw)
			fmt.Printf("Backstory: %s\n", npc.Backstory)
			fmt.Println("---------------------\n")
		},
	}
	npcCmd.AddCommand(generateCmd)
}
