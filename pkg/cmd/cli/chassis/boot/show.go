/*

 MIT License

 (C) Copyright 2023-2024 Hewlett Packard Enterprise Development LP

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
	"io"
	"strings"

	"github.com/stmcginnis/gofish/redfish"

	"github.com/Cray-HPE/gru/internal/query"

	"github.com/spf13/cobra"

	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
)

// NewShowCommand creates the `boot` subcommand for `show`.
func NewShowCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "boot host [...host]",
		Short: "Boot information",
		Long:  `Show the current BootOrder; BootNext, networkRetry, and more`,
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := query.Async(
				getBootInformation,
				hosts,
			)
			cli.PrettyPrint(content)
		},
	}
	return c
}

func getBootInformation(host string) interface{} {
	boot := Boot{Order: []string{}}
	c, err := auth.Connection(host)
	if err != nil {
		boot.Error = err
		return boot
	}

	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil || len(systems) < 1 {
		boot.Error = err
		return boot
	}

	bo := fmt.Sprintf(
		"%s/%s",
		strings.TrimRight(
			systems[0].ODataID,
			"/",
		),
		"BootOptions",
	)
	client := systems[0].GetClient()
	resp, err := client.Get(bo)

	// GigaByte has this key
	if err == nil || resp != nil {

		opts := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&opts)
		if err != nil {
			boot.Error = err
			return boot
		}

		for _, b := range systems[0].Boot.BootOrder {
			ep := fmt.Sprintf(
				"%s/%s",
				bo,
				strings.TrimPrefix(
					b,
					"Boot",
				),
			)
			resp, err := client.Get(ep)
			if err != nil {
				boot.Error = err
				return boot
			}
			names := make(map[string]interface{})
			err = json.NewDecoder(resp.Body).Decode(&names)
			if err != nil {
				boot.Error = err
				return boot
			}

			boot.Order = append(
				boot.Order,
				strings.TrimSpace(names["Description"].(string)),
			)
		}
	} else {

		bo = strings.TrimRight(
			systems[0].ODataID,
			"/",
		)
		response, err := client.Get(bo)
		if err != nil {
			boot.Error = err
			return boot
		}

		names := make(map[string]interface{})
		body, err := io.ReadAll(response.Body)
		err = json.Unmarshal(
			body,
			&names,
		)
		if err != nil {
			boot.Error = err
			return boot
		}

		bootMap := redfish.Boot{}
		bootJSON, err := json.Marshal(names["Boot"])
		if err != nil {
			boot.Error = err
			return boot
		}

		err = bootMap.UnmarshalJSON(bootJSON)
		if err != nil {
			boot.Error = err
			return boot
		}

		for _, v := range bootMap.BootOrder {
			boot.Order = append(
				boot.Order,
				strings.TrimSpace(v),
			)
		}

	}

	boot.Next = systems[0].Boot.BootNext

	return boot
}
