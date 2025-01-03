package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
	"github.com/spf13/cobra"
)

var port string

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires one argument: 'remote' or 'local'")
	}
	if args[0] != "remote" && args[0] != "local" {
		return fmt.Errorf("invalid argument: '%s'. Must be 'remote' or 'local'", args[0])
	}
	return nil
}

func generate(mode string) {
	switch mode {
	case "remote":
		cfg, err := internal.GenerateRemote(port)
		if err != nil {
			os.Exit(1)
		}
		b, err := json.Marshal(cfg)
		if err != nil {
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(b))
	case "local":
		cfg, err := internal.GenerateLocal(port)
		if err != nil {
			os.Exit(1)
		}
		b, err := json.Marshal(cfg)
		if err != nil {
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(b))
	default:
		// this is impossible
		os.Exit(1)
	}
}

var cfgCmd = &cobra.Command{
	Use:   "gen-config [remote|local]",
	Short: "generate either a local or remote config based on local machine's ip address",
	Long:  "generate either a local or remote config based on local machine's ip address. need to add tokens manually",
	Args:  validateArgs,
	Run: func(cmd *cobra.Command, args []string) {
		mode := args[0]
		generate(mode)
	},
}

func init() {
	cfgCmd.Flags().StringVarP(&port, "port", "p", internal.DEFAULT_PORT, "the port to listen/send data on")
	rootCmd.AddCommand(cfgCmd)
}
