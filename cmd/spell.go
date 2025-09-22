package cmd

import (
	"fmt"
	"strings"

	"dnd-cli/internal/data"

	"github.com/spf13/cobra"
)

// spellCmd represents the spell command
var spellCmd = &cobra.Command{
	Use:   "spell [spell name]",
	Short: "Looks up details for a D&D spell",
	Long: `Provides detailed information about a specified D&D spell, 
including its description, properties, and source. 

Examples:
  dnd spell "Fireball"
  dnd spell "eldritch blast"`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		spellName := strings.Join(args, " ")

		spell, err := data.GetSpellByName(spellName)
		if err != nil {
			fmt.Printf("Hark! Thy query, good sir or madam, doth bewilder my arcane senses. Pray tell, couldst thou rephrase thy plea, for its meaning doth elude my understanding: %v\n", err)
			return
		}

		fmt.Printf("\n--- %s ---\n", spell.Name)
		fmt.Printf("Description: %s\n", spell.Description)
		for key, value := range spell.Properties {
			fmt.Printf("%s: %v\n", key, value)
		}
		fmt.Printf("Source: %s (%s)\n", spell.Book, spell.Publisher)
		fmt.Println("-------------------\n")
	},
}

func init() {
	RootCmd.AddCommand(spellCmd)
}