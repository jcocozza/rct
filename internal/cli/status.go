package cli

import (
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "attempt to dial local server",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		server := internal.NewServer(cfg.Server.Addr, cfg.Server.Token, nil)
		alive := server.IsAlive()
		if alive {
			fmt.Fprintf(os.Stdout, "running on %s\n", server.Addr)
			return
		}
		fmt.Fprintf(os.Stdout, "server not running on %s\n", server.Addr)
		pid := internal.ServerExists()
		if pid != -1 {
			fmt.Fprintf(os.Stdout, "a server may be running on a different address (have you changed your .rct.json file?).\ncheck pid: %d\n", pid)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
