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
	go s.Run()
	for t := range results {
		if listenVerbose {
			v := drawTextBox(t)
			fmt.Fprintf(os.Stdout, "received: \n%s\n", v)
		}
	}
}
