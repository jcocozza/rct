package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
)

type command struct {
	Name string
	Usage string
}

var commands = map[string]command{
	"listen": {
		Name: "listen",
		Usage: "start a listener",
	},
	"kill": {
		Name: "kill",
		Usage: "kill the listener",
	},
	//"send": {
	//	Name: "send",
	//	Usage: "send text to server",
	//},
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <text>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "	or: %s <command> [OPTIONS]\n", os.Args[0])
	fmt.Fprint(os.Stderr, "Commands:\n")
	for _, cmd := range commands {
		fmt.Fprintf(os.Stderr, "	%s	%s\n", cmd.Name, cmd.Usage)
	}
	flag.PrintDefaults()
}

func Execute() {
	flag.Usage = usage
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
	case "kill":
		parseKill()
		runKill()
	default:
		txt := parseSend()
		runSend(cfg, txt)
	}
}
