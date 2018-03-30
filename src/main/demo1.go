package main

import (
	"fmt"
	"libbeat/common"
	"strings"
)

const commandLenBuffer = 50

var redisCommands = map[string]struct{}{
	"GET":              {},
	"ZCOUNT":           {},
	"ZUNIONSTORE":      {},
}

func main() {
	var key = common.NetString("get")
	var keys string = "get"


	var buf [commandLenBuffer]byte

	upper := buf[:len(key)]
	for i, b := range key {
		if 'a' <= b && b <= 'z' {
			b = b - 'a' + 'A'
		}
		upper[i] = b
		fmt.Println(i)
		fmt.Println(b)
	}
	fmt.Printf("%s \n",string(upper))
	_, exists := redisCommands[string(upper)]
	fmt.Println(exists)
	fmt.Print(strings.ToUpper(keys))
}
