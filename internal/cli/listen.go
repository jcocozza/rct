package cli

import (
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
	"github.com/spf13/cobra"
)

var verbose bool
var detach bool
var results chan string

func runListen(cfg internal.RCTConfig) {
	if verbose {
		results = make(chan string, 1)
	}
	s := internal.NewServer(cfg.Server.Addr, cfg.Server.Token, results)
	if verbose {
		fmt.Fprintf(os.Stdout, "listening on %s\n", s.Addr)
	}
	// run the server as detached
	if detach {
		pid, err := s.RunDetached()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to start listening: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%d\n", pid)
		os.Exit(0)
		return
	}
	// in verbose, we spawn the server and then block
	// the server will push text to the results chan for printing
	go func() {
		err := s.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to start listening %s\n", err.Error())
		}
	}()
	for t := range results {
		if verbose {
			v := drawTextBox(t)
			fmt.Fprintf(os.Stdout, "received: \n%s\n", v)
		}
	}
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "start listening",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		runListen(cfg)
	},
}

func init() {
	listenCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbosity")
	listenCmd.Flags().BoolVarP(&detach, "detach", "d", false, "detach the session and listen in background")
	listenCmd.MarkFlagsMutuallyExclusive("verbose", "detach")
	rootCmd.AddCommand(listenCmd)
}
