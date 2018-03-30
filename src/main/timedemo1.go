package main

import (
	"time"
	"fmt"
	"strconv"
)

func main() {
	t := time.Now()
	fmt.Println(t)

	fmt.Println(t.UTC().Format(time.UnixDate))
	fmt.Println(t.Unix())
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	fmt.Println(timestamp)
	timestamp = timestamp[:10]
	fmt.Println(timestamp)

	timestamp1 := time.Now().UnixNano() / 1000000
	fmt.Println(timestamp1)
	tt := time.Time{}
	fmt.Println(tt)

}