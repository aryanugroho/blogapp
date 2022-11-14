package console

import (
	"context"
	"os"
	"os/signal"
	"runtime"

	"github.com/aryanugroho/blogapp/app"
	"github.com/aryanugroho/blogapp/config"
	"github.com/aryanugroho/blogapp/internal/logger"
	"github.com/aryanugroho/blogapp/transport/mq/subscriber"

	"gopkg.in/ukautz/clif.v1"
)

func (c *Console) StartListener() *clif.Command {
	return clif.NewCommand("start", "Starting consumer server.", func(o *clif.Command, in clif.Input, out clif.Output) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			go callOnInterrupt(cancel)
		}()

		logger.Info(ctx, "runtime go version "+runtime.Version())

		// init app
		app, err := app.NewApp(ctx)
		if err != nil {
			return err
		}

		// init subscriber
		if config.All().Consumers.Enable {
			subscriber.Init(ctx, app)
		}

		logger.Info(ctx, "fraud listener is disabled")

		// Keep the program running until an interrupt signal is received.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		return nil
	})
}
