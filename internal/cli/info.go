package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "print current rct config",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		cfgBytes, err := json.Marshal(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading config: %s", err.Error())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(cfgBytes))
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
