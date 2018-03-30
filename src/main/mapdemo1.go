package main

import (
	"fmt"
	"sync"
	"bytes"
)

type userInfo struct {
	Name string
	Age  int
}
var m sync.Map

func main() {
	//返回键的现有值(如果存在)，否则存储并返回给定的值，如果是读取则返回true，如果是存储返回false。
	vv, ok := m.LoadOrStore("1", "one")
	fmt.Println(vv, ok) //one false

	vv, ok = m.LoadOrStore("1", "twasdo")
	fmt.Println(vv, ok)

	vv, ok = m.Load("1")
	fmt.Println(vv, ok) //one true

	vv, ok = m.LoadOrStore("1", "oneone")
	fmt.Println(vv, ok) //one true

	vv, ok = m.Load("1")
	fmt.Println(vv, ok) //one true

	m.Store("1", "oneone")
	vv, ok = m.Load("1")
	fmt.Println(vv, ok) // oneone true
//存储一个设置的键值
	m.Store("2", "two")
	m.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		return true
	})

	m.Delete("1")
	m.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		return true
	})

	map1 := make(map[string]userInfo)
	var user1 userInfo
	user1.Name = "ChamPly"
	user1.Age = 24
	map1["user1"] = user1

	var user2 userInfo
	user2.Name = "Tom"
	user2.Age = 18
	m.Store("map_test", map1)

	mapValue, _ := m.Load("map_test")

	for k, v := range mapValue.(interface{}).(map[string]userInfo) {
		fmt.Println(k, v)
		fmt.Println("name:", v.Name)
		fmt.Println("name:", v.Age + 1)
	}
	var buf bytes.Buffer
	buf.WriteString("1234" + ":")
	buf.WriteString("get" + ":")
	buf.WriteString("caochun_key")
	buf.String()
	fmt.Println(buf.String())
}
