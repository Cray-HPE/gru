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
	"fmt"
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stmcginnis/gofish/redfish"
)

// Boot represents boot configuration on the BMC. Only Error is emitted on empty.
type Boot struct {
	Order []string `json:"order,omitempty" yaml:"order,omitempty"`
	Next  string   `json:"next,omitempty" yaml:"next,omitempty"`
	Error error    `json:"error,omitempty" yaml:"error,omitempty"`
}

// Override represents the result of the boot override.
type Override struct {
	Target redfish.BootSourceOverrideTarget `json:"target" yaml:"target"`
	Error  error                            `json:"error,omitempty" yaml:"error,omitempty"`
}

// NewCommand creates the `boot` subcommand for `chassis`.
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:    "boot",
		Short:  "Set next boot device",
		Long:   `Overrides the next boot device for a one-time override. Only sets UEFI boot modes`,
		Hidden: false,
	}
	c.AddCommand(
		NewBiosOverrideCommand(),
		NewHddOverrideCommand(),
		NewPxeOverrideCommand(),
		NewUEFIHttpOverrideCommand(),
		NewNoneOverrideCommand(),
	)

	c.PersistentFlags().BoolP(
		"persist",
		"p",
		false,
		fmt.Sprintln(
			"Override continuously instead of onetime; persistent override",
		),
	)

	c.PersistentFlags().BoolP(
		"now",
		"n",
		false,
		fmt.Sprintln(
			"Reset the server(s) immediately after applying the boot override",
		),
	)

	return c
}

// issueOverride issues a boot override action against a host.
func issueOverride(host string, override interface{}) interface{} {
	o := Override{}
	v := viper.GetViper()

	c, err := auth.Connection(host)
	if err != nil {
		o.Error = err
		return o
	}

	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil || len(systems) < 1 {
		o.Error = err
		return o
	}

	boot := redfish.Boot{
		BootSourceOverrideTarget: override.(redfish.BootSourceOverrideTarget),
		BootSourceOverrideMode:   redfish.UEFIBootSourceOverrideMode,
	}

	if v.GetBool("persist") {
		boot.BootSourceOverrideEnabled = redfish.ContinuousBootSourceOverrideEnabled
	} else {
		boot.BootSourceOverrideEnabled = redfish.OnceBootSourceOverrideEnabled
	}

	err = systems[0].SetBoot(boot)
	o.Target = override.(redfish.BootSourceOverrideTarget)
	if err != nil {
		o.Error = err
	}

	return o
}
