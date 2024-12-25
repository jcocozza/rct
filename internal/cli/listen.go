package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
)

var (
	listenVerbose bool
	detach bool
	results chan string
)

func parseListen() {
	listen := flag.NewFlagSet("listen", flag.ExitOnError)
	listen.BoolVar(&listenVerbose, "v", false, "verbosity")
	listen.BoolVar(&detach, "d", false, "detach the session and listen in the background")

	// args[0] = app name
	// args[1] = "listen"
	err := listen.Parse(os.Args[2:])
	if err != nil {
		listen.Usage()
		return
	}
	if listenVerbose && detach {
		fmt.Fprintf(os.Stderr, "error: verbose and detach cannot both be set to true\n")
		os.Exit(1)
	}
	if listenVerbose {
		results = make(chan string, 1)
	}
}

func runListen(cfg internal.RCTConfig) {
	s := internal.NewServer(cfg.Server.Addr, cfg.Server.Token, results)
	if listenVerbose {
		fmt.Fprintf(os.Stdout, "listening on %s\n", s.Addr)
	}
	// run the server as detached
	if detach {
		pid, err := s.RunDetached()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to start listening: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "rct listening on: %d\n", pid)
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
		if listenVerbose {
			v := drawTextBox(t)
			fmt.Fprintf(os.Stdout, "received: \n%s\n", v)
		}
	}
}
