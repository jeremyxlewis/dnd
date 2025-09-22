package cmd

import (
	"dnd-cli/internal/tui"

	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launches the interactive TUI application",
	Long: `The 'tui' command starts the full-screen, interactive
text-based user interface for the D&D CLI companion.

Use this to access various features in a more interactive way.`,
	Run: func(cmd *cobra.Command, args []string) {
		tui.StartTUI()
	},
}

func init() {
	RootCmd.AddCommand(tuiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tuiCmd.PersistentFlags().String("foo", "bar", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tuiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
