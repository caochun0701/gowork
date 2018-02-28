package main

import (
	"packetbeat/protos"
	"sync"
	"fmt"
)

/*
主要用來过滤保存解析后的redis协议
使用RWMutex，读操作不会被锁，写操作保持同步
1、功能是找出热key
2、找出大values
*/

//定义一个结构体
type hotKye struct {
	Name string
	Age int
	//解析后的結果
	results protos.Reporter
}
var m sync.Map
//方法用来hotkey map操作
func hotKeyMatch(hotKye *hotKye){
	value,ok := m.LoadOrStore(hotKye.Name,hotKye.Age)
	//one false
	fmt.Println(value, ok)
	value1 ,ok := m.Load(hotKye.Name)
	if ok {
		m.Store(hotKye.Name,hotKye.Age +1)
		fmt.Println(value1)
	}
	value2 ,ok := m.Load(hotKye.Name)
	if ok {
		fmt.Println(value2)
	}
}
func main() {
	hotKye := &hotKye{
		Name:"caochun",
		Age:20,
	}
	hotKeyMatch(hotKye)
}
