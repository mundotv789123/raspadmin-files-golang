package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func Run() {
	cmd := &cli.Command{
		Name:    "Raspadmin",
		Usage:   "Raspadmin file manager web",
		Version: "0.1.0",
		Commands: []*cli.Command{
			&webCommand,
		},
		DefaultCommand: webCommand.Name,
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
