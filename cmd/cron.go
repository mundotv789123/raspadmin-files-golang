package cmd

import (
	"context"

	"github.com/mundotv789123/raspadmin/internal/database"
	icongenerator "github.com/mundotv789123/raspadmin/jobs/icon_generator"
	"github.com/urfave/cli/v3"
)

var cronCommand = cli.Command{
	Name:   "cron",
	Usage:  "Start the icon generator background service",
	Action: runCron,
}

func runCron(_ context.Context, cmd *cli.Command) error {
	_, err := database.OpenDbConnection()
	if err != nil {
		return err
	}
	icongenerator.StartBackgroundService()
	select {}
}

