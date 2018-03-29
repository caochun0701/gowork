package redis

import (
	"libbeat/beat"
	"sync"
	"libbeat/common"
	"fmt"
)

/*
1、功能是找出热key
2、找出大values
*/

/*
定义一个结构体
 */
type hotKey struct {
	event beat.Event
}

var m sync.Map

/*
 存储所有接收到Fields
*/
func suspectedHotkeyStore(event common.MapStr){
	port, ok := event["port"]
	if ok {
		fmt.Println(port)
	}

	method, ok := event["method"]
	if ok {
		fmt.Println(method)
	}

	resource, ok := event["resource"]
	if ok {
		fmt.Println(resource)
	}

	vv, ok := m.LoadOrStore("port:keyName", "one")
	if ok{
		//当前 Fields count +1
	}else{
		//当前Fields count = 1
	}
	debugf("%s , %v", vv, ok)
}
/*
 存储符合条件的大value Fields
*/
func bigValuesStore(event common.MapStr){

}



func main() {
	//格式化解析结果为map
	fields := common.MapStr{
		"type": "redis",
		"query": "key",
	}
	debugf("%s", fields)
}

