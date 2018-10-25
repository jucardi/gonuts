package cli

import (
	"github.com/jucardi/go-logger-lib/log"
	"github.com/jucardi/gonuts/deps"
	"github.com/jucardi/gonuts/workspace"
	"github.com/spf13/cobra"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manages workspaces for dependencies",
	//Run:   workspaceRun,
}

var wsSetCmd = &cobra.Command{
	Use: "set",
	Short: `Sets a workspace for the current project using the git IDs associated to the 'gonuts.yml' dependencies.
If no dependencies checkpoint exists with a 'gonuts.yml', creates a checkpoint of the current project first`,
	Run: wsSetRun,
}

var wsResetCmd = &cobra.Command{
	Use:   "reset",
	Short: `Leaves the workspace restoring the revisions to their original revision`,
	Run:   wsResetRun,
}

var wsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: `Shows the current workspace (if any)`,
	Run:   wsStatusRun,
}

func init() {
	workspaceCmd.AddCommand(wsSetCmd)
	workspaceCmd.AddCommand(wsResetCmd)
	workspaceCmd.AddCommand(wsStatusCmd)
}

func wsSetRun(cmd *cobra.Command, args []string) {
	log.Info("Setting workspace...")
	_, err := deps.Manager().Load()
	if err != nil && err != deps.ErrNoDepsExists {
		log.Fatal(err.Error())
	} else if err == deps.ErrNoDepsExists {
		depsSave()
	}

	if err := workspace.Manager().Set(); err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Done.")
}

func wsResetRun(cmd *cobra.Command, args []string) {
	log.Info("Resetting workspace...")
	if err := workspace.Manager().Reset(); err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Done")
}

func wsStatusRun(cmd *cobra.Command, args []string) {
	if ws, err := workspace.Manager().Load(); err != nil {
		log.Fatal(err.Error())
	} else if ws != nil {
		//log.Info(ws.Name)
		println(ws.Name)
	}
}
