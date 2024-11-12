package subscriber

import (
	"context"
	"log"
	"sync"

	"github.com/aryanugroho//config"
	"github.com/aryanugroho//internal/infrastructure/mq"
	"github.com/aryanugroho//internal/usecase"
)

type Subscriber struct {
	app *usecase.App
}

type Consumer struct {
	topic                  string
	subscription           string
	maxOutstandingMessages int
	numGoroutines          int
	handler                mq.Handler
	toggle                 bool
}

var s *Subscriber

func Init(ctx context.Context, app *usecase.App) error {
	s = &Subscriber{
		app: app,
	}

	client, err := mq.NewMQClient(ctx, config.All().GCloud.ProjectID)
	if err != nil {
		log.Fatalf("failed init mq pubsub with error %v", err)
		return err
	}

	consumers := []Consumer{}

	wg := new(sync.WaitGroup)
	for _, c := range consumers {
		if !c.toggle {
			continue
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, c Consumer) {
			defer wg.Done()
			st := c.handler
			client.Subscribe(c.topic, c.subscription, c.maxOutstandingMessages, c.numGoroutines, st)
		}(wg, c)
	}

	wg.Wait()
	return nil
}
