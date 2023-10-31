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

package boot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/Cray-HPE/gru/pkg/query"
	"github.com/spf13/cobra"
)

// NewShowCommand creates the `boot` subcommand for `show`.
func NewShowCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "boot [flags] host [...host]",
		Short: "Boot information",
		Long:  `Show the current BootOrder; BootNext, networkRetry, and more`,
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := query.Async(getBootInformation, hosts)
			cli.MapPrint(content)
		},
	}
	return c
}

func getBootInformation(host string, args ...string) interface{} {
	boot := cli.Boot{}
	c, err := auth.Connection(host)
	if err != nil {
		boot.Error = err
		return boot
	}

	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil {
		boot.Error = err
	}

	// BootOptions in gofish is unexported, so get it here manually
	// could also modify the vendored dir to expose it, making it easier here
	// rather than modify a stable external package, do some work here to get the friendly names
	bo := fmt.Sprintf("%s/%s", strings.TrimRight(systems[0].ODataID, "/"), "BootOptions")

	// get the bootoptions endpoint
	resp, err := systems[0].Client.Get(bo)
	if err != nil {
		return err
	}

	// store options in a map for easy manipulation
	opts := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&opts)
	if err != nil {
		return err
	}

	// make a map for the descriptions
	// boot.Descriptions = make(map[string]string, 0)
	for _, b := range systems[0].Boot.BootOrder {
		// the endpoint is BootOptions/NNNN so strip off 'Boot' from the boot order name
		ep := fmt.Sprintf("%s/%s", bo, strings.TrimPrefix(b, "Boot"))
		// get the endpoint
		response, err := systems[0].Client.Get(ep)
		if err != nil {
			return err
		}
		// decode to a map
		names := make(map[string]interface{})
		err = json.NewDecoder(response.Body).Decode(&names)
		if err != nil {
			return err
		}

		bd := cli.BootDescription{}
		bd[b] = names["Description"].(string)
		// create a key with the boot option using the friendly name as the value
		boot.Order = append(boot.Order, bd)
	}

	boot.Next = systems[0].Boot.BootNext
	// boot.Order = systems[0].Boot.BootOrder

	return boot
}
