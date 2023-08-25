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
	"github.com/stmcginnis/gofish/redfish"
)

// NewPowerOnCommand creates the `on` subcommand for `power`.
func NewPowerOnCommand() *cobra.Command {
	c := &cobra.Command{
		Use:    "on",
		Short:  "Power on",
		Long:   `Power on a node.`,
		Run:    powerOn,
		Hidden: false,
	}
	return c
}

// powerOn represents the cobra command that powers on nodes
func powerOn(cmd *cobra.Command, args []string) {
	action.Send(args, setPowerOn)
}

// setPowerOn gets the power state for a given host
func setPowerOn(host string) error {
	c, err := auth.Connection(host)
	if err != nil {
		return err
	}
	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil {
		return err
	}

	var resetType redfish.ResetType = redfish.OnResetType
	if powerOnCmd.Flags().Changed("force") {
		resetType = redfish.ForceOnResetType
	}

	err = systems[0].Reset(resetType)

	if err != nil {
		return err
	}

	return nil
}
