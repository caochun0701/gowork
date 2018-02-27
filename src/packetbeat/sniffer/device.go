package sniffer

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/tsg/gopacket/pcap"

	"github.com/elastic/beats/libbeat/logp"
)

var deviceAnySupported = runtime.GOOS == "linux"

// ListDevicesNames returns the list of adapters available for sniffing on
// this computer. If the withDescription parameter is set to true, a human
// readable version of the adapter name is added. If the withIP parameter
// is set to true, IP address of the adapter is added.
func ListDeviceNames(withDescription bool, withIP bool) ([]string, error) {
	//获取当前服务器所有的设备
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return []string{}, err
	}

	ret := []string{}
	for _, dev := range devices {
		//当前设备名字
		r := dev.Name
		//设备描述
		if withDescription {
			desc := "No description available"
			if len(dev.Description) > 0 {
				desc = dev.Description
			}
			r += fmt.Sprintf(" (%s)", desc)
		}
		//找出带IP的设备
		if withIP {
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
			}
			r += fmt.Sprintf(" (%s)", ips)

		}
		ret = append(ret, r)
	}
	return ret, nil
}
/*
解析设备名称
 */
func resolveDeviceName(name string) (string, error) {
	if name == "" {
		return "any", nil
	}
	//string到int
	if index, err := strconv.Atoi(name); err == nil {
		// Device is numeric id
		devices, err := ListDeviceNames(false, false)
		if err != nil {
			return "", fmt.Errorf("Error getting devices list: %v", err)
		}

		name, err = deviceNameFromIndex(index, devices)
		if err != nil {
			return "", fmt.Errorf("Couldn't understand device index %d: %v", index, err)
		}

		logp.Info("Resolved device index %d to device: %s", index, name)
	}

	return name, nil
}

func deviceNameFromIndex(index int, devices []string) (string, error) {
	if index >= len(devices) {
		return "", fmt.Errorf("Looking for device index %d, but there are only %d devices",
			index, len(devices))
	}

	return devices[index], nil
}
