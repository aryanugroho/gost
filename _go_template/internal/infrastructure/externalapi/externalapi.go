package externalapi

import (
	"context"

	"github.com/aryanugroho//config"
	"github.com/aryanugroho//internal/infrastructure/externalapi/api"
)

type ExternalAPI interface {
	Sample() api.SampleProvider
}

type ExternalAPIProvider struct {
	sample api.SampleProvider
}

func NewExternalAPI(ctx context.Context, externalAPIConfig config.ExternalAPIConfiguration) (ExternalAPI, error) {

	// external api provider initialization
	//httpClient3s := &http.Client{Timeout: 3 * time.Second}
	sampleProvider := api.NewSample()

	return &ExternalAPIProvider{
		sample: sampleProvider,
	}, nil
}

func (p *ExternalAPIProvider) Sample() api.SampleProvider {
	return p.sample
}
