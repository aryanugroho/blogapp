package main

import (
	"context"
	"os"

	"github.com/aryanugroho/blogapp/internal/config"
	"github.com/aryanugroho/blogapp/internal/logger"
	"github.com/aryanugroho/blogapp/internal/statsd"
	"github.com/aryanugroho/blogapp/transport/console"
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

	cli := clif.New("FraudService", "1.0.0", "Fraud service managed by risk team")
	cmd, err := console.Init()
	if err != nil {
		logger.Fatal(ctx, "failed init console", err)
	}
	cli.Add(cmd.StartServer())
	cli.Add(cmd.MigrateCreate())
	cli.Add(cmd.MigrateRun(ctx))
	cli.Add(cmd.MigrateRollback())
	cli.Add(cmd.MigrateReset())
	cli.Add(cmd.MigrateRefresh())

	cli.Run()
}
