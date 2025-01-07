package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jcocozza/rct/internal"
	"github.com/spf13/cobra"
)

const version string = "0.0.1"

// initialized when cli is called
var cfg internal.RCTConfig

func handleStdin() string {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading from stdin: %s", err.Error())
		os.Exit(1)
	}
	return strings.TrimSpace(string(b))
}

func runSend(c internal.RCTConfig, txt string) {
	if len(c.Delivery) == 0 {
		fmt.Fprintf(os.Stderr, "no remotes specified. check config file.\n")
		os.Exit(1)
	}
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
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stat, _ := os.Stdin.Stat()
		useStdin := (stat.Mode() & os.ModeCharDevice) == 0
		var txt string
		if useStdin {
			txt = handleStdin()
		} else {
			if len(args) == 0 {
				fmt.Fprintf(os.Stderr, "error: expected text (one argument)\n")
				os.Exit(1)
			}
			txt = args[0]
		}
		initConfig()
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
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
