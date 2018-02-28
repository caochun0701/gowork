package memqueue

import (
	"libbeat/logp"
)

type logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
}

var defaultLogger logger = logp.NewLogger("memqueue")
