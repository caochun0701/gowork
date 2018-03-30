package main

import (
	"libbeat/common"
	"time"
	//"packetbeat/protos/redis"
	"libbeat/beat"
	"fmt"
)

func main() {
	//格式化解析结果为map
	fields := common.MapStr{
		"type":         "redis",
		"status":       "ok",
		"responsetime": time.Now(),
		"method":       common.NetString("GET"),
		"resource":     "caochun_key",
		"query":        "get caochun_key",
		"bytes_in":     uint64(120),
		"bytes_out":    uint64(200),
		"port": 6403,
	}
	beat := beat.Event{
		Fields:fields,
	}
	fmt.Println(beat.Fields["port"])

	//redis.SuspectedHotkeyStore(beat.Fields)

}
