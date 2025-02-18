package cli

import (
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use: "ping",
	Short: "check if there is a connection to the server",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		if len(cfg.Delivery) == 0 {
			fmt.Fprintf(os.Stderr, "no remotes specified. check config file.\n")
			os.Exit(1)
		}
		for _, host := range cfg.Delivery {
			client := internal.NewClient(host.Addr, host.Token)
			err := client.Ping()
			if err != nil {
				fmt.Fprintf(os.Stdout, "%s is down\n", client.ServerAddr)
			} else {
				fmt.Fprintf(os.Stdout, "%s is up\n", client.ServerAddr)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
