package cli

import (
	"fmt"

	"github.com/jucardi/gonuts/cmd/gonuts/version"
	"github.com/spf13/cobra"
)

const (
	goListCmd = `go list -f '{{ if .Imports }}{{ join .Imports "\n" }}{{ end }}' ./... | xargs -L1 go list -f '{{ if not .Standard }}{{ .ImportPath  }}{{ end }}'`
	usage     = ``
	long      = ``
)

var rootCmd = &cobra.Command{
	Use:              "gonuts",
	Short:            "Golang package manager",
	Long:             fmt.Sprintf(long, version.Version, version.Built),
	PersistentPreRun: initCmd,
	Run:              run,
}

// Execute starts the execution of the parse command.
func Execute() {
	rootCmd.Flags().StringP("file", "f", "", "INPUT: A JSON or YAML file to use as an input for the data to be parsed")
	rootCmd.Flags().StringArrayP("definition", "d", []string{}, "Other templates to be loaded to be used in the 'templates' directive.")

	rootCmd.Execute()
}

func printUsage(cmd *cobra.Command) {
	cmd.Println(fmt.Sprintf(long, version.Version, version.Built))
	cmd.Usage()
}

func initCmd(cmd *cobra.Command, args []string) {
	// FromCommand(cmd)
	cmd.Use = fmt.Sprintf(usage, cmd.Use)
}

func run(cmd *cobra.Command, args []string) {
	// input, _ := cmd.Flags().GetString("file")
	// str, _ := cmd.Flags().GetString("string")
	// url, _ := cmd.Flags().GetString("url")
	// output, _ := cmd.Flags().GetString("output")
	// definitions, _ := cmd.Flags().GetStringArray("definition")
	// pattern, _ := cmd.Flags().GetString("pattern")
}
