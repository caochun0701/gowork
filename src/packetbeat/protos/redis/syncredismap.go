package redis

import (
	"libbeat/beat"
	"sync"
	"libbeat/common"
	"fmt"
	"bytes"
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
type CountNum struct {
	Count int
}

var m sync.Map

/*
 存储所有接收到Fields
*/
func SuspectedHotkeyStore(event common.MapStr){

	port, ok := event["port"].(*common.Endpoint)
	if ok {
		fmt.Println(port)
	}

	method, ok := event["method"].(*common.Endpoint)
	if ok {
		fmt.Println(method)
	}

	resource, ok := event["resource"].(*common.Endpoint)
	if ok {
		fmt.Println(resource)
	}
	//进行key拼接
	var buf bytes.Buffer
	buf.WriteString(string(port.Port) + ":")
	buf.WriteString(method.Method + ":")
	buf.WriteString(method.Resource)

	fields, ok := m.Load(buf.String())
	if ok{
		//当前 Fields count +1
		for key, value := range fields.(interface{}).(map[string]*common.Endpoint) {
			event["count"] = value.Count + 1
			m.Store(key,event)
		}
	}else{
		//当前Fields count = 1
		event["count"] = 1
		m.Store(buf,event)
	}
	debugf("%s , %v", fields, ok)
}
/*
 存储符合条件的大value Fields
*/
func bigValuesStore(event common.MapStr){

}

