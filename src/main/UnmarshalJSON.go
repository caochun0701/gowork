package main

import (
	"encoding/json"
	"fmt"
)

/*解析json对象*/

type Message struct {
	Asks []Order `json:"Bids"`
	Bids []Order `json:"Asks"`
}

type Order struct {
	Price  float64
	Volume float64
}

func (o *Order) UnmarshalJSON(data []byte) error{
	var v [2]float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.Price = v[0]
	o.Volume = v[1]
	return nil
}

func main() {
	b := []byte(`{"Asks": [[21, 1], [22, 1]] ,"Bids": [[20, 1], [19, 1]]}`)
	var m Message
	if err := json.Unmarshal(b,&m);err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", m)

}
