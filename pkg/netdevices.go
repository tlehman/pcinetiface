package netdevices

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

type Node struct {
	Class       string `json:"class"`
	Id          string `json:"id"`
	BusInfo     string `json:"businfo"`
	LogicalName string `json:"logicalname"`
	Config      struct {
		Children []Node `json:"children"`
	} `json:"configuration"`
}

type DeviceInfo struct {
	PCIAddress       string
	NetworkInterface string
	LinkState        string // "UP", "DOWN"
}

func GetPCINetInterfaces() ([]DeviceInfo, error) {
	cmd := exec.Command("lshw", "-class", "network", "-json")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var lshwData []Node
	err = json.Unmarshal(out, &lshwData)
	if err != nil {
		return nil, err
	}

	var devices []DeviceInfo
	findPCIDevices(lshwData, &devices)
	return devices, nil
}

func findPCIDevices(nodes []Node, devices *[]DeviceInfo) {
	for _, node := range nodes {
		if node.Class == "network" {
			pciAddress := ""
			networkInterface := ""
			if strings.HasPrefix(node.BusInfo, "pci") {
				pciAddress = strings.TrimPrefix(node.BusInfo, "pci:")
				networkInterface = node.LogicalName
				linkState, err := getLinkState(networkInterface)

				if err != nil {
					fmt.Printf("error getting link state: %s\n", err)
				}
				*devices = append(*devices, DeviceInfo{
					PCIAddress:       pciAddress,
					NetworkInterface: networkInterface,
					LinkState:        linkState,
				})
			}
		}
	}
}

func getLinkState(networkInterface string) (string, error) {
	linkStatePath := fmt.Sprintf("/sys/class/net/%s/operstate", networkInterface)
	content, err := ioutil.ReadFile(linkStatePath)
	if err != nil {
		return "", err
	}

	linkState := strings.TrimSpace(string(content))
	return linkState, nil
}
