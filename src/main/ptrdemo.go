package main

import "fmt"

const MAX  int = 4


/*go 指针使用*/
func main() {

	/* 声明实际变量 */
	var a int= 20

	/* 声明指针变量 */
	var b *int

	/* 指针变量的存储地址 */
	b = &a

	fmt.Printf("a 变量的地址是: %x\n", &a  )

	/* 指针变量的存储地址 */
	fmt.Printf("b 变量储存的指针地址: %x\n", b )

	/* 使用指针访问值 */
	fmt.Printf("*b 变量的值: %d\n", *b )

	/**
	指针数组
	 */
 	c := []int {1,2,34,5}

	for i := 0;i< MAX ;i ++  {
		fmt.Printf("index: %d value is %d \n", i,c[i])
	}
	/**
	整数型指针数组
	 */
	var ptr [MAX]*int
	for i :=0; i < MAX ;i++  {
		/* 整数地址赋值给指针数组 */
		ptr[i] = &c[i]
	}

	for  i := 0; i < MAX; i++ {
		fmt.Printf("a[%d] = %d \n",i,*ptr[i])
	}

}
