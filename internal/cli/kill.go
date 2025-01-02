package cli

import (
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
	"github.com/spf13/cobra"
)

func runKill() {
	pid, err := internal.Kill()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to kill: %s\n", err.Error())
		os.Exit(1)
		return
	}
	fmt.Fprintf(os.Stdout, "killed %d\n", pid)
}

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill the background listener if it is running",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runKill()
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
