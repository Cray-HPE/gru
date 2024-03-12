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
	"fmt"
	"github.com/Cray-HPE/gru/internal/set"
	"github.com/Cray-HPE/gru/pkg/cmd"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stmcginnis/gofish/redfish"
)

// NewPowerOffCommand creates the `off` subcommand for `power`.
func NewPowerOffCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "off host [...host]",
		Short: "Power off the target machine(s)",
		Long: `Powers off the target machine(s) with an ACPI shutdown.
Permits forcing a shutdown (without waiting for the OS),
as well as a power-button emulated shutdown.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			var resetType redfish.ResetType

			hosts := cli.ParseHosts(args)

			v := viper.GetViper()
			bindErr := v.BindPFlags(c.Flags())
			cmd.CheckError(bindErr)

			resetType = redfish.GracefulShutdownResetType
			if v.GetBool("force") {
				resetType = redfish.ForceOffResetType
			}
			if v.GetBool("button") {
				resetType = redfish.PushPowerButtonResetType
			}

			content := set.Async(Issue, hosts, resetType)
			cli.PrettyPrint(content)
		},
	}
	c.PersistentFlags().BoolP(
		"force",
		"f",
		false,
		fmt.Sprintln(
			"Immediately power off without waiting for the OS",
		),
	)
	c.PersistentFlags().BoolP(
		"button",
		"b",
		false,
		"Emulate a power-button press",
	)
	c.MarkFlagsMutuallyExclusive("button", "force")
	return c
}
