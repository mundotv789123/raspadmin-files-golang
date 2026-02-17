package cmd

import (
	ctx "context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/database"
	icongenerator "github.com/mundotv789123/raspadmin/jobs/icon_generator"
	"github.com/mundotv789123/raspadmin/router"
	"github.com/robfig/cron"
	"github.com/urfave/cli/v3"
)

var webCommand = cli.Command{
	Name:   "web",
	Usage:  "Start the web server",
	Action: runWeb,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "gen-thumbnail",
			Usage: "Enable cron to generate file thumbnails",
		},
	},
}

var cronIsRunning = false

func runWeb(_ ctx.Context, cmd *cli.Command) error {
	_, err := database.OpenDbConnection()
	if err != nil {
		return err
	}

	genThumb := cmd.Bool("gen-thumbnail")
	if genThumb {
		c := cron.New()
		c.AddFunc("0 1 * * * *", func() {
			defer func() { cronIsRunning = false }()
			if !cronIsRunning {
				cronIsRunning = true
				icongenerator.RunGenerator()
			} else {
				slog.Warn("cron is running, ignoring...")
			}
		})
		c.Start()
	}

	r := gin.Default()

	webCtx := &router.WebContext{DB: database.DB}
	webCtx.Routers(r)

	return r.Run("0.0.0.0:8080")
}
