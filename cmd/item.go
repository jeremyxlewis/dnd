package cmd

import (
	"fmt"
	"strings"

	"dnd-cli/internal/data"

	"github.com/spf13/cobra"
)

// itemCmd represents the item command
var itemCmd = &cobra.Command{
	Use:   "item [item name]",
	Short: "Looks up details for a D&D item",
	Long: `Provides detailed information about a specified D&D item, 
including its name and a brief description. 

Examples:
  dnd item "Potion of Healing"
  dnd item "Longsword"`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		itemName := strings.Join(args, " ")

		item, err := data.GetItemByName(itemName)
		if err != nil {
			fmt.Printf("Hark! Thy query, good sir or madam, doth bewilder my arcane senses. Pray tell, couldst thou rephrase thy plea, for its meaning doth elude my understanding: %v\n", err)
			return
		}

		fmt.Printf("\n--- %s ---\n", item.Name)
		fmt.Printf("Description: %s\n", item.Description)
		fmt.Println("-------------------\n")
	},
}

func init() {
	RootCmd.AddCommand(itemCmd)
}