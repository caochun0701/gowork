package main

import "fmt"

func main() {
	/*指向指针的指针*/
	var a int = 2000
	var ptr *int
	var pptr **int
	/* 指针 ptr 地址 */
	ptr = &a
	/* 指向指针 ptr 地址 */
	pptr = &ptr
	/* 获取 pptr 的值 */
	fmt.Printf("变量 a = %d\n",a)
	fmt.Printf("指针变量 *ptr = %d\n",*ptr)
	fmt.Printf("指向指针的指针变量 **pptr = %d\n",**pptr)

}
