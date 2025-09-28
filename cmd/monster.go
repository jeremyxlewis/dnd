package cmd

import (
	"fmt"
	"strings"

	"dnd-cli/internal/data"

	"github.com/spf13/cobra"
)

// monsterCmd represents the monster command
var monsterCmd = &cobra.Command{
	Use:   "monster [monster name]",
	Short: "Looks up details for a D&D monster",
	Long: `Provides detailed information about a specified D&D monster, 
including its name and a brief description. 

Examples:
  dnd monster "Goblin"
  dnd monster "Ancient Red Dragon"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		monsterName := strings.Join(args, " ")

		monster, err := data.GetMonsterByName(monsterName)
		if err != nil {
			fmt.Printf("Hark! Thy query, good sir or madam, doth bewilder my arcane senses. Pray tell, couldst thou rephrase thy plea, for its meaning doth elude my understanding: %v\n", err)
			return
		}

		fmt.Printf("\n--- %s ---\n", monster.Name)
		fmt.Printf("Description: %s\n", monster.Description)
		fmt.Print("-------------------\n")
	},
}

func init() {
	RootCmd.AddCommand(monsterCmd)
}
