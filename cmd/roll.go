package cmd

import (
	"fmt"
	"strings"

	"dnd-cli/internal/dice"

	"github.com/spf13/cobra"
)

var (
	advantage    bool
	disadvantage bool
)

// rollCmd represents the roll command
var rollCmd = &cobra.Command{
	Use:   "roll [notation]",
	Short: "Rolls dice using standard D&D notation (e.g., 2d6, 1d20+5)",
	Long: `Rolls dice based on the provided notation. 

Examples:
  dnd roll 1d20
  dnd roll 2d6+3
  dnd roll 1d20 --advantage
  dnd roll 4d6 --disadvantage`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		notation := strings.ToLower(args[0])

		if advantage && disadvantage {
			fmt.Println("Hark! Thou canst not have both advantage and disadvantage, for the fates are fickle but not contradictory!")
			return
		}

		dr, err := dice.ParseDiceNotation(notation)
		if err != nil {
			fmt.Printf("Hark! Thy query, good sir or madam, doth bewilder my arcane senses. Pray tell, couldst thou rephrase thy plea, for its meaning doth elude my understanding: %v\n", err)
			return
		}

		var total int
		var rolls []int

		if advantage {
			roll1, _ := dr.Roll()
			roll2, _ := dr.Roll()
			total = max(roll1, roll2)
			fmt.Printf("Rolling %s with Advantage: (Roll 1: %d, Roll 2: %d) -> Total: %d\n", dr.Notation, roll1, roll2, total)
		} else if disadvantage {
			roll1, _ := dr.Roll()
			roll2, _ := dr.Roll()
			total = min(roll1, roll2)
			fmt.Printf("Rolling %s with Disadvantage: (Roll 1: %d, Roll 2: %d) -> Total: %d\n", dr.Notation, roll1, roll2, total)
		} else {
			total, rolls = dr.Roll()
			fmt.Printf("Rolling %s: %v -> Total: %d\n", dr.Notation, rolls, total)
		}
	},
}

func init() {
	RootCmd.AddCommand(rollCmd)

	rollCmd.Flags().BoolVarP(&advantage, "advantage", "a", false, "Roll with advantage (roll twice, take higher)")
	rollCmd.Flags().BoolVarP(&disadvantage, "disadvantage", "d", false, "Roll with disadvantage (roll twice, take lower)")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}