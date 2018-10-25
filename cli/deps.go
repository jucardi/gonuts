package cli

import (
	"fmt"
	"github.com/jucardi/go-logger-lib/log"
	"github.com/jucardi/gonuts/deps"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Creates or updates a dependencies snapshot",
	Run:   depsRun,
}

func init() {
	depsCmd.Flags().BoolP("dry-run", "d", false, "Shows the dependencies info in the console instead of saving the snapshot")
}

func depsRun(cmd *cobra.Command, args []string) {
	dryrun, _ := cmd.Flags().GetBool("dry-run")
	if dryrun {
		depsDryRun()
	} else {
		depsSave()
	}
}

func depsSave() (*deps.Dependencies) {
	log.Info("Gathering dependencies information...")
	dependencies, err := deps.Manager().Generate()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Saving dependencies information...")
	if err := dependencies.Save(); err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Done.")
	return dependencies
}

func depsDryRun() {
	dInfo, err := deps.Manager().Generate()
	if err != nil {
		log.Fatal(err.Error())
	}

	println()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprintln(w, "Package Name\tRevision")
	for k, v := range dInfo.Dependencies {
		fmt.Fprintf(w, " %s\t %s", k, v.Revision)
	}
	w.Flush()
}
