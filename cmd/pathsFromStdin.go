package cmd

import (
	"github.com/gkwa/fewport/core"
	"github.com/spf13/cobra"
)

var pathsFromStdinCmd = &cobra.Command{
	Use:     "paths-from-stdin",
	Aliases: []string{"pfs"},
	Short:   "Clean up Google URLs in markdown files specified via stdin",
	Long:    `This command reads file paths from stdin and cleans up Google URLs in the specified markdown files by removing specific query parameters.`,
	Run: func(cmd *cobra.Command, args []string) {
		core.ProcessPathsFromStdin(cmd.Context(), core.CleanGoogleURLsInFile)
	},
}

func init() {
	rootCmd.AddCommand(pathsFromStdinCmd)
}
