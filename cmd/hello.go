package cmd

import (
	"github.com/gkwa/fewport/core"
	"github.com/spf13/cobra"
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "files-from-dir",
	Aliases: []string{"ffd"},
	Short: "Clean up Google URLs in markdown files",
	Long:  `This command recursively searches for markdown files in the specified directory and cleans up Google URLs by removing specific query parameters.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Usage()
			return
		}

		dir := args[0]
		err := core.CleanGoogleURLs(dir)
		if err != nil {
			cmd.PrintErrf("Failed to clean Google URLs: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
}
