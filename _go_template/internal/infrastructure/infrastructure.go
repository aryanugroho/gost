package infrastructure

import (
	"context"

	"github.com/aryanugroho//config"
	"github.com/aryanugroho//internal/infrastructure/externalapi"
	"github.com/aryanugroho//internal/infrastructure/mq"
	"github.com/aryanugroho//internal/infrastructure/sqlstore"

	"github.com/getsentry/sentry-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Infrastructure is the wrapper for infra dependencies.
type Infrastructure interface {
	SQLStore() sqlstore.Store
	ExternalAPI() externalapi.ExternalAPI
	MQ() mq.Client
}

type Infra struct {
	sqlStore    sqlstore.Store
	externalAPI externalapi.ExternalAPI
	mq          mq.Client
}

func NewInfra(ctx context.Context, config config.Configuration) (Infrastructure, error) {
	// init sql store
	sqlStore, err := sqlstore.NewSQLStore(ctx, config.Database.Master, config.Database.Slave)
	if err != nil {
		return nil, err
	}

	// init external api
	externalAPI, err := externalapi.NewExternalAPI(ctx, config.ExternalAPI)
	if err != nil {
		return nil, err
	}

	// init mq
	mqClient, err := mq.NewMQClient(ctx, config.GCloud.ProjectID)
	if err != nil {
		return nil, err
	}

	// start datadog apm
	if config.App.ENV != "local" {

		tracer.Start()
		defer func() {
			tracer.Stop()
		}()
	}

	// init sentry
	if config.Sentry.DSN != nil {
		if err = sentry.Init(sentry.ClientOptions{
			Dsn:              *config.Sentry.DSN,
			Release:          config.App.Version,
			AttachStacktrace: true,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				if &event.Request != nil {
					event.Request.Cookies = ""
					event.Request.QueryString = ""
					event.Request.Headers = nil
					event.Request.Data = ""
				}
				return event
			},
			//TracesSampler: sentry.NewRateLimitSampler(1),
		}); err != nil {
			return nil, err
		}
	}

	// init others
	return &Infra{
		sqlStore:    sqlStore,
		externalAPI: externalAPI,
		mq:          mqClient,
	}, nil
}

func (i *Infra) SQLStore() sqlstore.Store {
	return i.sqlStore
}

func (i *Infra) ExternalAPI() externalapi.ExternalAPI {
	return i.externalAPI
}

func (i *Infra) MQ() mq.Client {
	return i.mq
}
