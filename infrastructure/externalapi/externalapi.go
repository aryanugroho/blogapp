package externalapi

import (
	"context"

	"github.com/aryanugroho/blogapp/config"
	"github.com/aryanugroho/blogapp/infrastructure/externalapi/sample"
)

type ExternalAPI interface {
	Sample() sample.SampleProvider
}

type ExternalAPIProvider struct {
	sample sample.SampleProvider
}

func NewExternalAPI(ctx context.Context, externalAPIConfig config.ExternalAPIConfiguration) (ExternalAPI, error) {

	// external api provider initialization
	//httpClient3s := &http.Client{Timeout: 3 * time.Second}
	sampleProvider := sample.NewSample()

	return &ExternalAPIProvider{
		sample: sampleProvider,
	}, nil
}

func (p *ExternalAPIProvider) Sample() sample.SampleProvider {
	return p.sample
}
