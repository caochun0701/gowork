package codec

import "libbeat/beat"

type Codec interface {
	Encode(index string, event *beat.Event) ([]byte, error)
}
