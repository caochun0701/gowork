package main

import (
	"fmt"
)

func main() {
	//定义一个数组
	numbers := [6]int {1,2,3,4,5,6}
	//for循环
	for _,value:=range numbers{

		fmt.Printf("value is :%d \n",value)
	}
	//for循环比较方式
	var a int
	var b int = 6
	for a < b{
		fmt.Printf("a %d < b  %d\n", a,b)
		a ++
	}
	//输出带有索引的
	for index,value :=range numbers  {
		fmt.Printf("第 %d 位 x 的值 = %d\n", index,value)
	}

}
