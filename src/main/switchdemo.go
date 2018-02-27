package main

import "fmt"

func main() {
	/* 定义局部变量 */
	var grade string = "B"

	var masks int = 90

	switch masks {
		case 90:
			grade = "A"
		case 80:
			grade = "B"
		case 50:
			grade = "C"
		default:
			grade = "D"
	}
	switch {
		case grade == "A":
			fmt.Printf("优秀!\n" )
		case grade == "B",grade == "C":
			fmt.Printf("良好\n" )
		default:
			fmt.Printf("及格\n" )
	}
	/*判断某个 interface 变量中实际存储的变量类型*/

	var x interface{}

	switch i := x.(type){
		case nil:
			fmt.Printf(" x 的类型 :%T",i)
		case int:
			fmt.Printf("x 是 int 型")
		case float64:
			fmt.Printf("x 是 float64 型")
		case func(int) float64:
			fmt.Printf("x 是 func(int) 型")
		case bool, string:
			fmt.Printf("x 是 bool 或 string 型" )
		default:
			fmt.Printf("未知型")
	}

}
