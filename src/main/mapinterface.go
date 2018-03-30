package main

import "log"

func init(){
	log.Println("this is init")
}

func main() {

	list := map[string]interface{}{
		"name":"caochun",
		"age":30,
	}

	for key,value := range list{
		log.Println(key)
		log.Println(value)
	}
	e,ok := list["name"]
	if ok {
		log.Println(e)
	}

}
