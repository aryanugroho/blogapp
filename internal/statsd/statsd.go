package statsd

import (
	dogStatsD "github.com/DataDog/datadog-go/v5/statsd"
	"github.com/aryanugroho/blogapp/internal/config"
)

var statsDClient dogStatsD.ClientInterface

const sampleRate = 0

func Init() {
	client, err := dogStatsD.New(config.GetStatsDAddress())
	if err != nil {
		statsDClient = &dogStatsD.NoOpClient{}
		return
	}
	statsDClient = client
}

func getStatsDClient() dogStatsD.ClientInterface {
	if statsDClient == nil {
		Init()
	}
	return statsDClient
}

func Incr(name string) error {
	return getStatsDClient().Incr(name, []string{}, sampleRate)
}

func Gauge(name string, value float64) error {
	return getStatsDClient().Gauge(name, value, []string{}, sampleRate)
}

func Close() error {
	return getStatsDClient().Close()
}
