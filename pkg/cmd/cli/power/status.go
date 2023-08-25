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
	"github.com/Cray-HPE/gru/pkg/action"
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/spf13/cobra"
)

// NewPowerStatusCommand creates the `status` subcommand for `power`.
func NewPowerStatusCommand() *cobra.Command {
	c := &cobra.Command{
		Use:    "status",
		Short:  "Power status",
		Long:   `Check power status.`,
		Run:    powerStatus,
		Hidden: false,
	}
	return c
}

// powerStatus represents the cobra command that queries the power status of nodes
func powerStatus(cmd *cobra.Command, args []string) {
	action.Get(args, getPowerStatus)
}

// getPowerStatus gets the power state for a given host
func getPowerStatus(host string) (string, error) {
	c, err := auth.Connection(host)
	if err != nil {
		return "", err
	}
	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil {
		return "", err
	}

	return string(systems[0].PowerState), nil
}
