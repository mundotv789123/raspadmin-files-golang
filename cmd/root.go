package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

func Run() {
	cmd := &cli.Command{
		Name:    "Raspadmin",
		Usage:   "Raspadmin file manager web",
		Version: "0.1.0",
		Commands: []*cli.Command{
			&cronCommand,
			&webCommand,
		},
		DefaultCommand: webCommand.Name,
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error(fmt.Sprintf("%s", err))
	}
}
