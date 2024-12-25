package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
)

func parseKill() {
	kill := flag.NewFlagSet("kill", flag.ExitOnError)
	kill.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s kill\n", os.Args[0])
		fmt.Fprint(os.Stderr, "  kill the background listener\n")
		kill.PrintDefaults()
	}
	// args[0] = app name
	// args[1] = "kill"
	err := kill.Parse(os.Args[2:])
	if err != nil {
		kill.Usage()
		return
	}
}

func runKill() {
	pid, err := internal.Kill()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to kill: %s\n", err.Error())
		os.Exit(1)
		return
	}
	fmt.Fprintf(os.Stdout, "killed %d\n", pid)

}
