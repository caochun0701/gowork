package includes

import (
	// load supported output plugins
	_ "libbeat/outputs/elasticsearch"
	_ "libbeat/outputs/fileout"
	//_ "libbeat/outputs/kafka"
	//_ "libbeat/outputs/logstash"
	// load support output codec
	_ "libbeat/outputs/codec/format"
	_ "libbeat/outputs/codec/json"
)
