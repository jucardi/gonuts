package cli

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetches the project dependencies into the project workspace environment",
	Run:   getRun,
}

func getRun(cmd *cobra.Command, args []string) {
	//c := exec.Command("go", "get -u ./...")
}
