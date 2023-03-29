package main

import (
	"fmt"

	nd "github.com/tlehman/pcinetiface/pkg/netdevices"
)

func main() {
	devices, err := nd.GetPCINetInterfaces()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, device := range devices {
		fmt.Printf(
			"PCIAddress: %s, NetworkInterface: %s, LinkState: %s\n",
			device.PCIAddress,
			device.NetworkInterface,
			device.LinkState,
		)
	}
}
