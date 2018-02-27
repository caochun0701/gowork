package main

import "fmt"

/*结构体作为函数参数*/
type Bookss struct{

	subject string

	title string

	author string

	book_id int

}

func printBook(book Bookss){
	fmt.Printf( "Book title : %s\n", book.title);
	fmt.Printf( "Book author : %s\n", book.author);
	fmt.Printf( "Book subject : %s\n", book.subject);
	fmt.Printf( "Book book_id : %d\n", book.book_id);
}

func printBookPtr(book *Bookss){
	fmt.Printf( "Book title : %s\n", book.title);
	fmt.Printf( "Book author : %s\n", book.author);
	fmt.Printf( "Book subject : %s\n", book.subject);
	fmt.Printf( "Book book_id : %d\n", book.book_id);
}


func main() {
	var Book1 Bookss        /* 声明 Book1 为 Books 类型 */
	var Book2 Bookss        /* 声明 Book2 为 Books 类型 */
	/* book 1 描述 */
	Book1.title = "Go 语言"
	Book1.author = "www.runoob.com"
	Book1.subject = "Go 语言教程"
	Book1.book_id = 6495407
	/* book 2 描述 */
	Book2.title = "Python 教程"
	Book2.author = "www.runoob.com"
	Book2.subject = "Python 语言教程"
	Book2.book_id = 6495700
	/* 打印 Book1 信息 */
	printBook(Book1)
	/* 打印 Book2 信息 */
	printBook(Book2)

	/*结构指针体*/

	/* 打印 Book1 信息 */
	printBookPtr(&Book1)
	/* 打印 Book2 信息 */
	printBookPtr(&Book2)




}
