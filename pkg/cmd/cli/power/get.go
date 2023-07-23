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

package power

import (
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/Cray-HPE/gru/pkg/query"
	"github.com/spf13/cobra"
)

// NewGetCommand creates the `power` subcommand for `get`.
func NewGetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "power",
		Short: "List available power actions",
		Long:  `Gets the available power actions for a host.`,
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := query.Async(getPowerActionInformation, hosts)
			cli.MapPrint(content)
		},
	}
	return c
}

func getPowerActionInformation(host string) interface{} {
	action := Action{}
	c, err := auth.Connection(host)
	if err != nil {
		action.Error = err
		return action
	}

	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil {
		action.Error = err
	}

	// FIXME: Gigabyte does not return available power commands.
	action.Actions = systems[0].SupportedResetTypes
	return action
}
