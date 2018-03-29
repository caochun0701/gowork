package main

import (
	"fmt"
	"reflect"
	"unsafe"
	"github.com/go-lumber/log"
)

/*切片*/
func printSlice(x []int){
	fmt.Printf("len=%d cap=%d slice=%v\n",len(x),cap(x),x)
}

func main() {
	var numbers  = make([]int,3,5)
	printSlice(numbers)
	numbers = append(numbers, 100)
	fmt.Print(numbers)

	s := []string{"w","s","k"}
	d := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	log.Println(*d)

}
