package main

import (
	"fmt"
)

/**
Go 语言中 range 关键字用于for循环中迭代数组(array)、切片(slice)、
通道(channel)或集合(map)的元素。
在数组和切片中它返回元素的索引值，在集合中返回 key-value 对的 key 值。
*/

func main() {
	nums := []int {1,2,3,4}
	sum := 0
	for _,num := range nums{
		sum +=num
	}
	fmt.Println("sum:",sum)
	kvs := map[string]string {"a":"apple","b":"banana"}
	for key,value :=range kvs {
		fmt.Printf("%s -> %s\n", key, value)
	}

	//range也可以用来枚举Unicode字符串。第一个参数是字符的索引，第二个是字符（Unicode的值）本身。
	for i, c := range "go" {
		fmt.Println(i, c)
	}
}
