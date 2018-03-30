package main

import (
	"github.com/go-lumber/log"
	"time"
)

func main() {

	log.Println(time.Now().Format(time.ANSIC))
	time.After(time.Millisecond * 2000)
	log.Println(time.Time{})
	log.Println(time.Now().Format(time.ANSIC))
}
