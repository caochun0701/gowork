package elasticsearch

import (
	"libbeat/common"
	"libbeat/monitoring"
)

func makeSnapshot(R *monitoring.Registry) common.MapStr {
	mode := monitoring.Full
	return common.MapStr(monitoring.CollectStructSnapshot(R, mode, false))
}
