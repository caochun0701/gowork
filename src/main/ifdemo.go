package main

import "fmt"

/**
go 条件判断语句
 */
func main() {
	/*定义一个局部变量*/

	var a int = 5
	if a < 10{
		fmt.Printf("a value of %d \n",a)
		a++
		fmt.Printf("a value of %d \n",a)
	}
	//if else 语句
	if a < 10{
		fmt.Printf("a value of %d \n",a)
	}else{
		fmt.Printf("a value of %d \n",a)
	}

	if a < 10{
		if a ==6{
			fmt.Printf("a value of %d \n",a)
		}
	}

}
