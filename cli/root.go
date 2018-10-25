package cli

import (
	"fmt"
	"github.com/jucardi/go-logger-lib/log"
	"github.com/jucardi/gonuts/version"
	"github.com/spf13/cobra"
)

const (
	usage = ``
	long  = ``
)

var rootCmd = &cobra.Command{
	Use:              "gonuts",
	Short:            "Golang package manager",
	Long:             fmt.Sprintf(long, version.Version, version.Built),
	PersistentPreRun: initCmd,
}

// Execute starts the execution of the parse command.
func Execute() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enables Verbose mode.")
	rootCmd.AddCommand(depsCmd)
	rootCmd.AddCommand(workspaceCmd)
	rootCmd.Execute()
}

func printUsage(cmd *cobra.Command) {
	cmd.Println(fmt.Sprintf(long, version.Version, version.Built))
	cmd.Usage()
}

func initCmd(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	cmd.Use = fmt.Sprintf(usage, cmd.Use)
}
