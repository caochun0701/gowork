package main

import (
<<<<<<< HEAD
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
=======
	"github.com/go-lumber/log"
	"time"
)

func main() {

	log.Println(time.Now().Format(time.ANSIC))
	time.After(time.Millisecond * 2000)
	log.Println(time.Time{})
	log.Println(time.Now().Format(time.ANSIC))
}
>>>>>>> fa4227b0048ff97b13e35f7c61909be47d8dfea2
