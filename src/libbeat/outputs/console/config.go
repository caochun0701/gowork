package console

import "libbeat/outputs/codec"

type Config struct {
	Codec codec.Config `config:"codec"`

	// old pretty settings to use if no codec is configured
	Pretty bool `config:"pretty"`

	BatchSize int
}

var defaultConfig = Config{}
