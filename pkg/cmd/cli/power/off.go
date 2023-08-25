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

// NewPowerOffCommand creates the `off` subcommand for `power`.
func NewPowerOffCommand() *cobra.Command {
	c := &cobra.Command{
		Use:    "off",
		Short:  "Power off",
		Long:   `Power off a node.`,
		Run:    powerOff,
		Hidden: false,
	}
	return c
}

// powerOff represents the cobra command that powers off nodes
func powerOff(cmd *cobra.Command, args []string) {
	action.Send(args, setPowerOff)
}

// setPowerOff gets the power state for a given host
func setPowerOff(host string) error {
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

	var resetType redfish.ResetType = redfish.GracefulShutdownResetType
	if powerOffCmd.Flags().Changed("force") {
		resetType = redfish.ForceOffResetType
	}
	if powerOffCmd.Flags().Changed("button") {
		resetType = redfish.PushPowerButtonResetType
	}
	if powerOffCmd.Flags().Changed("restart") {
		resetType = redfish.GracefulRestartResetType
	}
	if powerOffCmd.Flags().Changed("nmi") {
		resetType = redfish.NmiResetType
	}
	if powerOffCmd.Flags().Changed("restart") && powerOffCmd.Flags().Changed("force") {
		resetType = redfish.ForceRestartResetType
	}

	err = systems[0].Reset(resetType)

	if err != nil {
		return err
	}

	return nil
}
