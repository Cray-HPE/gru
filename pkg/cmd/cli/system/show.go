/*

 MIT License

 (C) Copyright 2023 Hewlett Packard Enterprise Development LP

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

package system

import (
	"fmt"
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/spf13/cobra"
	"sync"
)

// NewShowCommand creates the `system` subcommand for `show`.
func NewShowCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "system [flags] host [...host]",
		Short: "System information.",
		Long: `
Show the Server Manufacturer, Server Model, System Version, and Firmware Version for the given server(s).
`,
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := query(hosts)
			cli.MapPrint(content)
		},
	}
	return c
}

// query displays some basic information about the node(s).
// Displays the server model name, System version, and BMC firmware version.
func query(hosts []string) map[string]interface{} {
	var wg sync.WaitGroup

	sliceLength := len(hosts)
	wg.Add(sliceLength)

	fmt.Printf("Querying BMCs for [%5d] nodes ... \n", len(hosts))

	sm := make(map[string]interface{})

	for _, host := range hosts {
		go func(host string) {
			defer wg.Done()
			sm[host] = getSystemInformation(host)
		}(host)
	}
	wg.Wait()
	return sm
}

// getSystemInformation gets the Manufacturer, Model, BIOSVersion, and FirmwareVersion of the host.
func getSystemInformation(host string) System {
	c := auth.Connection(host)
	defer c.Logout()

	// Retrieve the service root
	service := c.Service

	// FIXME: Return on failure, do not halt execution.
	// Query the systems data using the session token
	managers, err := service.Managers()
	cmd.CheckError(err)
	systems, err := service.Systems()
	cmd.CheckError(err)
	bios := System{
		Manufacturer:    systems[0].Manufacturer,
		Model:           systems[0].Model,
		BIOSVersion:     systems[0].BIOSVersion,
		FirmwareVersion: managers[0].FirmwareVersion,
	}
	return bios
}
