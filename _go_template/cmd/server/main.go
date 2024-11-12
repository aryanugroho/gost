package main

import (
	"context"
	"os"

	"github.com/aryanugroho//internal/presenter/console"
	"github.com/aryanugroho//pkg/config"
	"github.com/aryanugroho//pkg/logger"
	"github.com/aryanugroho//pkg/statsd"

	"gopkg.in/ukautz/clif.v1"
)

func main() {
	ctx := context.Background()
	_ = config.Load("./config.yaml")
	logger.Init()
	statsd.Init()

	// No need to run CLI if there is no argument
	if len(os.Args) == 1 {
		return
	}

	cli := clif.New("", "1.0.0", "")
	cmd, err := console.Init()
	if err != nil {
		logger.Fatal(ctx, "failed init console", err)
	}
	cli.Add(cmd.StartServer())
	cli.Add(cmd.MigrateCreate())
	cli.Add(cmd.MigrateRun(ctx))
	cli.Add(cmd.MigrateRollback())

	cli.Run()
}
