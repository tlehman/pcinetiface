package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"
)

type Lshw struct {
	XMLName xml.Name `xml:"node"`
	Nodes   []Node   `xml:"node"`
}

type Node struct {
	XMLName xml.Name `xml:"node"`
	Class   string   `xml:"class,attr"`
	Id      string   `xml:"id,attr"`
	Nodes   []Node   `xml:"node"`
}

type DeviceInfo struct {
	PCIAddress        string
	NetworkInterface  string
}

func main() {
	devices, err := getPCIDevices()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, device := range devices {
		fmt.Printf("PCIAddress: %s, NetworkInterface: %s\n", device.PCIAddress, device.NetworkInterface)
	}
}

func getPCIDevices() ([]DeviceInfo, error) {
	cmd := exec.Command("lshw", "-xml")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return nil, err
	}

	var lshwData Lshw
	err = xml.Unmarshal(out.Bytes(), &lshwData)
	if err != nil {
		return nil, err
	}

	var devices []DeviceInfo
	for _, node := range lshwData.Nodes {
		if node.Class == "network" {
			pciAddress := ""
			networkInterface := ""
			for _, childNode := range node.Nodes {
				if childNode.Class == "pci" {
					pciAddress = strings.TrimPrefix(childNode.Id, "pci:")
				} else if childNode.Class == "interface" {
					networkInterface = childNode.Id
				}
			}
			if pciAddress != "" && networkInterface != "" {
				devices = append(devices, DeviceInfo{PCIAddress: pciAddress, NetworkInterface: networkInterface})
			}
		}
	}

	return devices, nil
}

