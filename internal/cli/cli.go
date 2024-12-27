package cli

import (
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
	"github.com/spf13/cobra"
)

const version string = "0.0.1"

// initalized when cli is called
var cfg internal.RCTConfig

func runSend(c internal.RCTConfig, txt string) {
	if verbose {
		fmt.Fprintf(os.Stdout, "sending text:\n%s\n", drawTextBox(txt))
	}
	for _, host := range c.Delivery {
		client := internal.NewClient(host.Addr, host.Token)
		err := client.Send(txt)
		if err != nil && verbose {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		}
	}
}

var rootCmd = &cobra.Command{
	Use:     "rct [text to send]",
	Version: version,
	Short:   "rct is a tool for sending text between remote and local",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		txt := args[0]
		runSend(cfg, txt)
	},
}

func initConfig() {
	c, err := internal.ReadConfig()
	cobra.CheckErr(err)
	cfg = c
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbosity")
	cobra.OnInitialize(initConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
