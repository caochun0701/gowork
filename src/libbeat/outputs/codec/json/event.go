package json

import (
	"time"

	"libbeat/beat"
	"libbeat/common"
)

// Event describes the event structure for events
// (in-)directly send to logstash
type event struct {
	Timestamp time.Time     `struct:"@timestamp"`
	Meta      meta          `struct:"@metadata"`
	Fields    common.MapStr `struct:",inline"`
}

type eventfile struct {
	Timestamp time.Time     `struct:"@timestamp"`
	Fields    common.MapStr `struct:",inline"`
}

// Meta defines common event metadata to be stored in '@metadata'
type meta struct {
	Beat    string                 `struct:"beat"`
	Type    string                 `struct:"type"`
	Version string                 `struct:"version"`
	Fields  map[string]interface{} `struct:",inline"`
}

func makeEvent(index, version string, in *beat.Event) event {
	return event{
		Timestamp: in.Timestamp,
		Meta: meta{
			Beat:    index,
			Version: version,
			Type:    "doc",
			Fields:  in.Meta,
		},
		Fields: in.Fields,
	}
}

func makeEventFile(in *beat.Event) eventfile {
	return eventfile{
		Timestamp: in.Timestamp,
		Fields: in.Fields,
	}
}
