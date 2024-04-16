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

package power

import (
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/spf13/cobra"
	"github.com/stmcginnis/gofish/redfish"
)

// NewCommand creates the `power` subcommand for `chassis`.
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:    "power host [...host]",
		Short:  "Power Control",
		Long:   `Check power status, or power on, off, cycle, or reset a host`,
		Hidden: false,
	}
	c.AddCommand(
		NewPowerOffCommand(),
		NewPowerResetCommand(),
		NewPowerOnCommand(),
		NewPowerCycleCommand(),
		NewPowerStatusCommand(),
		NewPowerNMICommand(),
	)
	return c
}

// StateChange represents a change in power states.
type StateChange struct {
	PreviousPowerState  redfish.PowerState `json:"previousPowerState,omitempty" yaml:"previous_power_state,omitempty"`
	RequestedPowerState redfish.ResetType  `json:"requestedPowerState,omitempty" yaml:"requested_power_state,omitempty"`
	Error               error              `json:"error,omitempty" yaml:"error,omitempty"`
}

// State represents a single power state.
type State struct {
	PowerState redfish.PowerState `json:"powerState" yaml:"power_state"`
	Error      error              `json:"error,omitempty" yaml:"error,omitempty"`
}

// Issue issues an action against a host.
func Issue(host string, action interface{}) interface{} {
	sc := StateChange{}
	c, err := auth.Connection(host)
	if err != nil {
		sc.Error = err
		return sc
	}
	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil || len(systems) < 1 {
		sc.Error = err
		return sc
	}
	sc.PreviousPowerState = systems[0].PowerState
	sc.RequestedPowerState = action.(redfish.ResetType)
	err = systems[0].Reset(sc.RequestedPowerState)
	if err != nil {
		sc.Error = err
	}

	return sc
}
