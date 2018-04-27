package publish

import (
	"sync"
	"libbeat/beat"
	"libbeat/common"
	"packetbeat/config"
	"fmt"
	"bytes"
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
	key_name := event["key_name"].(common.NetString)
	//进行key拼接
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprint(port)+ ":")
	buf.WriteString(string(method) + ":")
	buf.WriteString(string(key_name))

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

func findHotKeysBigValues(m *sync.Map,client beat.Client,redisConf config.ProtocolCommon)  {
	m.Range(func(key, value interface{}) bool {
		//fmt.Println(value)
		fields := value.(common.MapStr)
		count := fields["count"].(int)
		bytesOut := fields["bytes_out"].(uint64)
		//找出热key
		if(count > redisConf.HotKeysCount){
			//返回给event进行输出
			client.Publish(beat.Event{Fields:fields})
		}
		//大values
		if(bytesOut > redisConf.BigValueSize){
			//返回给event进行输出
			client.Publish(beat.Event{Fields:fields})
		}

		return true
	})
}

