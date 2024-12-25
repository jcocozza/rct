package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: expected 'listen' or text to send\n")
		flag.Usage()
		return
	}
	cfg, err := internal.ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to read rct config: %s\n", err.Error())
		os.Exit(1)
	}

	switch os.Args[1] {
	case "listen":
		parseListen()
		runListen(cfg)
	default:
		txt := parseSend()
		runSend(cfg, txt)
	}
}
