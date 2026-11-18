/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	file string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "resume-yaml",
	Short: "Resume as YAML",
	Long: `
Keep your entire resume in one YAML file: experience, projects,
skills, and education, all in a single source. 

Tag each bullet with keywords like "go", "backend", or "leadership",
then tailor your resume for any application by passing keywords or a
job posting URL. resume-manager picks the most relevant content and
renders it as Markdown, plain text, or HTML.

Run it with no arguments to use the interactive TUI, or pass commands
and flags to skip it entirely for scripts and automation.

To get started, read the README:
https://github.com/jonathantrevino/resume-yaml#readme
`,
	Version: "0.1.0",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.resume-yaml.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "path to resume YAML file")

}
