package publish

import (
	"sync"
	"libbeat/beat"
	"libbeat/common"
	"fmt"
	"bytes"
	"time"
)

/*
1、功能是找出热key
2、找出大values
*/


/*
 存储所有接收到Fields
*/
func suspectedHotkeyStore(m *sync.Map, event common.MapStr){
	port := event["port"].(uint16)
	method := event["method"].(common.NetString)
	resource := event["resource"].(common.NetString)
	//进行key拼接
	var buf bytes.Buffer
	buf.WriteString(string(port)+ ":")
	buf.WriteString(string(method) + ":")
	buf.WriteString(string(resource))

	fields, ok := m.Load(buf.String())
	if ok{
		//当前 Fields count +1
		fs, ok := fields.(common.MapStr)
		if !ok {
			fmt.Println("!ok")
		}
		count := fs["count"]
		c, ok := count.(int)
		if !ok {
			fmt.Println("!ok")
		}
		fs["count"] = c + 1
		m.Store(buf.String(),fields)
	}else{
		//当前Fields count = 1
		event["count"] = 1
		m.Store(buf.String(),event)
	}
	debugf("%s , %v", fields, ok)
}

func findHotKeysBigValues(m *sync.Map,client beat.Client)  {
	m.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		//找出热key
		//if(){
			//返回给event进行输出
			//client.Publish(*pub)
		//}
		//大values
		//if(){
			//返回给event进行输出
			//client.Publish(*pub)
		//}

		return true
	})
}
