package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
	"github.com/spf13/cobra"
)

var cfgCmd = &cobra.Command{
	Use: "gen-config [server address] [delivery addresses]...",
	Short: "generate a config file for rct",
	Long: "generate a config file for rct. if you'd like a server token, you must add to the config file yourself.",
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "at least server address is required\n")
			os.Exit(1)
		}
		host := args[0]
		var deliveries []string
		if len(args) > 1 {
			deliveries = args[1:]
		}
		cfg, err := internal.GenerateConfig(host, deliveries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to generate config: %s\n", err.Error())
			os.Exit(1)
		}
		b, err := json.Marshal(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to marshal config: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(b))
	},
}

func init() {
	rootCmd.AddCommand(cfgCmd)
}
