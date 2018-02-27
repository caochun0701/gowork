package main

import (
	"encoding/json"
	"fmt"
)

/**
json 格式化
*/

type Student struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Guake bool `json:"guake"`
}


func main() {
	st := &Student {
		"Xiao Ming",
		16,
		true,
	}
	restJson,err := json.Marshal(st)
	if err != nil {
		return
	}
	fmt.Print(string(restJson))

}
