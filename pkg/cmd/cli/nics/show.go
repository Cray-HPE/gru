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

package nics

import (
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/Cray-HPE/gru/pkg/query"
	"github.com/spf13/cobra"
)

// NewShowCommand creates the `nics` subcommand for `show`.
func NewShowCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "nics [flags] host [...host]",
		Short: "Network interface information",
		Long:  `Show available network interface hardware`,
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := query.Async(getNICInformation, hosts)
			cli.MapPrint(content)
		},
		Hidden: true, // TODO: Remove or set to false once implemented.
	}
	return c
}

func getNICInformation(host string) interface{} {
	nics := NIC{}
	c, err := auth.Connection(host)
	if err != nil {
		nics.Error = err
		return nics
	}

	defer c.Logout()

	// service := c.Service

	if err != nil {
		nics.Error = err
	}
	return nics
}
