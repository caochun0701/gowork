package main

import (
	"net/http"
	"io/ioutil"
	"strings"
	"log"
)

func httpGet(){
	response, err := http.Get("http://www.01happy.com/demo/accept.php?id=1")
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(body))
}

func httpPost(){
	resp, err := http.Post("http://www.01happy.com/demo/accept.php",
		"application/x-www-form-urlencoded",
		strings.NewReader("name=cjb"))
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	log.Println(string(body))
	log.Println(resp.Status)
	log.Println(resp.StatusCode)
}


func main() {
	httpGet()
	httpPost()
}
