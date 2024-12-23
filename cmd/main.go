package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/jcocozza/rct/internal"
)

func deliver(hosts []internal.Host, txt string, verbose bool) bool {
	hasErr := false
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, h := range hosts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := internal.NewClient(h.Addr, h.Token)
			err := c.Send(txt)
			if err != nil {
				mu.Lock()
				hasErr = true
				mu.Unlock()
			}
			if err != nil && verbose {
				fmt.Fprintf(os.Stderr, "error: failed sending to %s: %s\n", h.Addr, err.Error())
			}
		}()
	}
	wg.Wait()
	return hasErr
}

func parse(cfg internal.RCTConfig) {
	verboseBase := flag.Bool("verbose", false, "enable verbosity")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <text>\n", os.Args[0])
		flag.PrintDefaults()
	}

	listen := flag.NewFlagSet("listen", flag.ExitOnError)
	verboseListen := listen.Bool("verbose", false, "enable verbosity")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: provide a command or argument\n")
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "listen":
		err := listen.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("error: failed to parse listen: %s\n", err.Error())
			os.Exit(1)
		}
		s := internal.NewServer(cfg.Server.Addr, cfg.Server.Token)
		if *verboseListen {
			fmt.Fprintf(os.Stdout, "listening on: %s\n", s.Addr)
		}
		s.Run()
	default:
		flag.Parse()
		if flag.NArg() != 1 {
			fmt.Fprintf(os.Stderr, "error: include text to send\n")
			flag.Usage()
			os.Exit(1)
		}
		args := flag.Args()
		txt := args[0]
		if *verboseBase {
			fmt.Fprintf(os.Stdout, "sending text: %s", txt)
		}
		b := deliver(cfg.Delivery, txt, *verboseBase)
		if b {
			fmt.Fprintf(os.Stderr, "error: failed to send\n")
			flag.Usage()
			os.Exit(1)
		}
	}
}

func main() {
	c, err := internal.ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read rct config: %s\n", err.Error())
		os.Exit(1)
	}
	parse(c)
}
