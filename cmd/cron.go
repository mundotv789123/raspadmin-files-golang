package cmd

import (
	"context"

	"github.com/mundotv789123/raspadmin/internal/database"
	icongenerator "github.com/mundotv789123/raspadmin/jobs/icon_generator"
	"github.com/urfave/cli/v3"
)

var cronCommand = cli.Command{
	Name:   "cron",
	Usage:  "Start the cron job",
	Action: runCron,
}

func runCron(_ context.Context, cmd *cli.Command) error {
	_, err := database.OpenDbConnection()
	if err != nil {
		return err
	}
	erro := icongenerator.RunGenerator()
	if erro != nil {
		return erro
	}
	return nil
}
