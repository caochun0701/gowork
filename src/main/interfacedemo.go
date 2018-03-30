package main

import "github.com/go-lumber/log"

//定义一个接口
type People interface{
	Name() string
}

type Student1 struct{
	name string
}

func (stu *Student1) Name() string{
	return stu.name
}

func getPeople() People {
	var stu *Student1
	return stu
}

//定义一个接口
type MyInterface interface {
	Print()
	Hello()
}

func main() {
	log.Println(getPeople())
	if getPeople() == nil {
		log.Println("AAA")
	}else{
		log.Println("BBB")
	}
	var st People = &Student1{"caochun"}
	log.Println(st.Name())

}
