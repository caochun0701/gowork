package main

import (
	"time"
	"log"
)

func main() {

	t1 := time.NewTimer(time.Second * 6)
	t2 := time.NewTimer(time.Second * 8)

	for{
		select{
		case <- t1.C:
			log.Println("t1")
			t1.Reset(time.Second * 5)
		case <- t2.C:
			log.Println("t2")
			t2.Reset(time.Second * 5)
		}
	}


}
