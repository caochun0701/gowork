package main

import (
	"github.com/tsg/gopacket/pcap"
	"log"
	"fmt"
	"strconv"
)

func main() {

	devices, err := pcap.FindAllDevs()
	if err != nil{
		log.Println("err",err)
	}

	log.Println(devices)
	for _, dev := range devices {


		ips := "Not assigned ip address"
		if len(dev.Addresses) > 0 {
			ips = ""

			for i, address := range []pcap.InterfaceAddress(dev.Addresses) {
				// Add a space between the IP address.
				if i > 0 {
					ips += " "
				}

				ips += fmt.Sprintf("%s", address.IP.String())
			}
			log.Println(ips)
		}
	}
	var name string = "any"

	index, err := strconv.Atoi(name)
	if err != nil {
		log.Print(err)
	}
	log.Println(index)
}
