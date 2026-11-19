/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"resume-yaml/internal/paths"

	"github.com/spf13/cobra"
)

// pathsCmd represents the paths command
var pathsCmd = &cobra.Command{
	Use:   "paths",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := paths.FindPaths(file)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("vault path: %v", p)
	},
}

func init() {
	rootCmd.AddCommand(pathsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pathsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pathsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
