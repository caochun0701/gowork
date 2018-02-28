package redis

import (
	"packetbeat/protos"
	"libbeat/beat"
	"sync"
	"fmt"
	"libbeat/common"
)

/*
主要用來过滤保存解析后的redis协议
使用RWMutex，读操作不会被锁，写操作保持同步
1、功能是找出热key
2、找出大values
*/

//定义一个结构体
type hotKye struct {
	//redis解析后的返回map
	event beat.Event
	//解析后的結果
	results protos.Reporter
}
var m sync.Map
//方法用来hotkey map操作
func hotKeyMatch(hotKye *hotKye){
	vv, ok := m.LoadOrStore("caochun", "one")
	//one false
	fmt.Println(vv, ok)
}
func main() {
	//格式化解析结果为map
	fields := common.MapStr{
		"type": "redis",
		"query": "key",
	}
	hotKye  := hotKye{
		event:beat.Event{Fields:fields},
	}
	hotKeyMatch(&hotKye)
}

