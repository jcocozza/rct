package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/jcocozza/rct/internal"
	"github.com/spf13/cobra"
)

func run(remoteIP string) {
	var username string
	fmt.Print("enter username: ")
	_, err := fmt.Scanln(&username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning: %s\n", err.Error())
		os.Exit(1)
	}
	var password string
	fmt.Print("enter password: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning: %s\n", err.Error())
		os.Exit(1)
	}

	err = internal.SendRemote(
		context.Background(),
		port,
		username,
		password,
		remoteIP,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to send: %s\n", err.Error())
		os.Exit(1)
	}
}

var distributeCmd = &cobra.Command{
	Use:   "distribute [remote ip addr list]",
	Short: "send a remote config to the specified ip addresses via scp.\nyou will be prompted for username/password",
	Long:  "send a remote config to the specified ip addresses via scp.\nyou will be prompted for username/password.\nrun 'gen-config remote' to see what config will be sent",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, ip := range args {
			run(ip)
		}
	},
}

func init() {
	distributeCmd.Flags().StringVarP(&port, "port", "p", internal.DEFAULT_PORT, "the port to listen/send data on")
	rootCmd.AddCommand(distributeCmd)
}
