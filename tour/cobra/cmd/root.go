package cmd

import "github.com/spf13/cobra"

var rootComd = &cobra.Command{}

func Execute() error {
	return rootComd.Execute()
}

func init() {
	rootComd.AddCommand(wordCmd)
	rootComd.AddCommand(timeCmd)
}
