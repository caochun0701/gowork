package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"libbeat/common"
	"libbeat/common/bus"
)

func TestGenerateHints(t *testing.T) {
	tests := []struct {
		event  bus.Event
		result bus.Event
	}{
		// Empty events should return empty hints
		{
			event:  bus.Event{},
			result: bus.Event{},
		},
		// Only kubernetes payload must return only kubernetes as part of the hint
		{
			event: bus.Event{
				"kubernetes": common.MapStr{
					"pod": common.MapStr{
						"name": "foobar",
					},
				},
			},
			result: bus.Event{
				"kubernetes": common.MapStr{
					"pod": common.MapStr{
						"name": "foobar",
					},
				},
			},
		},
		// Kubernetes payload with container info must be bubbled to top level
		{
			event: bus.Event{
				"kubernetes": common.MapStr{
					"container": common.MapStr{
						"name": "foobar",
						"id":   "abc",
					},
				},
			},
			result: bus.Event{
				"kubernetes": common.MapStr{
					"container": common.MapStr{
						"name": "foobar",
						"id":   "abc",
					},
				},
				"docker": common.MapStr{
					"container": common.MapStr{
						"name": "foobar",
						"id":   "abc",
					},
				},
			},
		},
		// Scenarios being tested:
		// logs/multiline.pattern must be a nested common.MapStr under hints.logs
		// metrics/module must be found in hints.metrics
		// not.to.include must not be part of hints
		// period is annotated at both container and pod level. Container level value must be in hints
		{
			event: bus.Event{
				"kubernetes": common.MapStr{
					"annotations": common.MapStr{
						"co.elastic.logs/multiline.pattern": "^test",
						"co.elastic.metrics/module":         "prometheus",
						"co.elastic.metrics/period":         "10s",
						"co.elastic.metrics.foobar/period":  "15s",
						"not.to.include":                    "true",
					},
					"container": common.MapStr{
						"name": "foobar",
						"id":   "abc",
					},
				},
			},
			result: bus.Event{
				"kubernetes": common.MapStr{
					"annotations": common.MapStr{
						"co.elastic.logs/multiline.pattern": "^test",
						"co.elastic.metrics/module":         "prometheus",
						"not.to.include":                    "true",
						"co.elastic.metrics/period":         "10s",
						"co.elastic.metrics.foobar/period":  "15s",
					},
					"container": common.MapStr{
						"name": "foobar",
						"id":   "abc",
					},
				},
				"hints": common.MapStr{
					"logs": common.MapStr{
						"multiline": common.MapStr{
							"pattern": "^test",
						},
					},
					"metrics": common.MapStr{
						"module": "prometheus",
						"period": "15s",
					},
				},
				"docker": common.MapStr{
					"container": common.MapStr{
						"name": "foobar",
						"id":   "abc",
					},
				},
			},
		},
	}

	cfg := defaultConfig()

	p := Provider{
		config: cfg,
	}
	for _, test := range tests {
		assert.Equal(t, p.generateHints(test.event), test.result)
	}
}
