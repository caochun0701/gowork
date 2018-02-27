package main

import (
	"os"
	"bufio"
	"fmt"
)

/*go 中defer操作
  在函数返回前执行一些操作
*/

func main() {
	//打开文件
	file, err := os.Open("/home/caochun/cc.log")
	//defer 关闭文件
	defer file.Close()

	if err != nil{
		return
	}
	//按行读取文件
	rd := bufio.NewReader(file)
	for{
		line, err := rd.ReadString('\n')
		if err != nil{
			break
		}
		fmt.Printf("%s ",line)
	}
}
