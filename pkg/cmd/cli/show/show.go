/*
MIT License

(C) Copyright 2022 Hewlett Packard Enterprise Development LP

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the "Software"),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

*/

package show

import (
	"container/list"
	"fmt"
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/spf13/cobra"
	"sync"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "show",
		Short: "Show information from a BMC.",
		Long:  `Show information from a BMC such as a BIOS briefing.`,
	}
	return c
}

type systemInfo struct {
	bios bios
	boot boot
}

type bios struct {
	manufacturer    string
	model           string
	biosVersion     string
	firmwareVersion string
}

type boot struct {
	next  string
	order list.List
}

// printNodeBriefing displays some basic information about the node(s).
// Displays the server model name, BIOS version, and BMC firmware version.
func printNodeBriefing(hosts []string) {
	var wg sync.WaitGroup

	sliceLength := len(hosts)
	wg.Add(sliceLength)

	fmt.Printf("Querying BMCs for [%5d] nodes ... \n", len(hosts))

	sm := make(map[string]interface{})

	for _, host := range hosts {
		go func(host string) {
			defer wg.Done()
			sysInfo := systemInfo{}
			sm[host] = getSystemInformation(host, sysInfo)
		}(host)
	}
	wg.Wait()
	fmt.Println("Finished")
}

// getSystemInformation gets the Manufacturer, Model, BIOSVersion, and FirmwareVersion of the host.
func getSystemInformation(host string, sysInfo systemInfo) systemInfo {
	c := auth.Connection(host)
	defer c.Logout()

	// Retrieve the service root
	service := c.Service

	// FIXME: Return on failure, do not halt execution.
	// Query the systems data using the session token
	managers, err := service.Managers()
	if err != nil {
		panic(err)
	}
	systems, err := service.Systems()
	if err != nil {
		panic(err)
	}
	bios := bios{
		manufacturer:    systems[0].Manufacturer,
		model:           systems[0].Model,
		biosVersion:     systems[0].BIOSVersion,
		firmwareVersion: managers[0].FirmwareVersion,
	}
	sysInfo.bios = bios
	return sysInfo
}
