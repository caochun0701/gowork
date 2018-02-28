package queuetest

import (
	"time"

	"libbeat/beat"
	"libbeat/common"
	"libbeat/publisher"
)

func makeEvent(fields common.MapStr) publisher.Event {
	return publisher.Event{
		Content: beat.Event{
			Timestamp: time.Now(),
			Fields:    fields,
		},
	}
}
