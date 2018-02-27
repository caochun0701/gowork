package main

import "fmt"

/**
map操作
*/

type User struct {
	Name string `json:"name"`
	Age int	`json:"age"`
}

func main() {
	user := map[int]User {1:User{"张三",20},2:User{"li", 9},3:User{"李四",23}}
	fmt.Println(user)
	for key,value := range user{
		fmt.Printf("key: %d value: %s",key,value)
	}
}
