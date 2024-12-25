package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
)

var (
	sendVerbose bool
)

func parseSend() string {
	send := flag.NewFlagSet("", flag.ExitOnError)
	send.BoolVar(&sendVerbose, "v", false, "verbosity")
	// args[0] = app name
	err := send.Parse(os.Args[1:])
	if err != nil {
		send.Usage()
		os.Exit(1)
	}
	if send.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "error: expects a single string of text")
		send.Usage()
		os.Exit(1)
	}
	txt := send.Args()[0]
	return txt
}

func runSend(cfg internal.RCTConfig, txt string) {
	if sendVerbose {
		fmt.Fprintf(os.Stdout, "sending text:\n%s\n", drawTextBox(txt))
	}
	for _, host := range cfg.Delivery {
		client := internal.NewClient(host.Addr, host.Token)
		err := client.Send(txt)
		if err != nil && sendVerbose {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		}
	}
}
