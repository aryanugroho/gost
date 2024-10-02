package console

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"[base_project]/internal/presenter/rest/api1"
	"[base_project]/internal/usecase"
	"[base_project]/pkg/logger"

	"gopkg.in/ukautz/clif.v1"
)

func (c *Console) StartServer() *clif.Command {
	return clif.NewCommand("start", "starting http server.", func(o *clif.Command, in clif.Input, out clif.Output) error {

		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			go callOnInterrupt(cancel)
		}()

		logger.Info(ctx, "runtime go version "+runtime.Version())

		// init app
		app, err := usecase.NewApp(ctx)
		if err != nil {
			return err
		}

		// init ws api
		api1.Init(ctx, app)

		return nil
	})
}

func callOnInterrupt(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
	cancel()
}
