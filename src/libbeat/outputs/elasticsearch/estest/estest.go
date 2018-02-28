package estest

import (
	"time"

	"libbeat/outputs/elasticsearch"
	"libbeat/outputs/elasticsearch/internal"
	"libbeat/outputs/outil"
)

// GetTestingElasticsearch creates a test client.
func GetTestingElasticsearch(t internal.TestLogger) *elasticsearch.Client {
	client, err := elasticsearch.NewClient(elasticsearch.ClientSettings{
		URL:              internal.GetURL(),
		Index:            outil.MakeSelector(),
		Username:         internal.GetUser(),
		Password:         internal.GetUser(),
		Timeout:          60 * time.Second,
		CompressionLevel: 3,
	}, nil)
	internal.InitClient(t, client, err)
	return client
}
