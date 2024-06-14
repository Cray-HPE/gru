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

// NewPowerCycleCommand creates the `cycle` subcommand for `power`.
func NewPowerCycleCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "cycle host [...host]",
		Short: "Power cycle the target machine(s)",
		Long: `Performs an ACPI shutdown and startup to power cycle the target machine(s).
Also allows bypassing the OS shutdown, forcing a warm boot.`,
		Run: func(c *cobra.Command, args []string) {
			var resetType redfish.ResetType

			hosts := cli.ParseHosts(args)

			v := viper.GetViper()
			bindErr := v.BindPFlags(c.Flags())
			cmd.CheckError(bindErr)

			resetType = redfish.GracefulRestartResetType
			if v.GetBool("force") {
				resetType = redfish.ForceRestartResetType
			}

			content := set.Async(Issue, hosts, resetType)
			cli.PrettyPrint(content)
		},
		Hidden: false,
	}
	c.PersistentFlags().BoolP(
		"force",
		"f",
		false,
		fmt.Sprintln(
			"Immediately restart waiting for the OS (warm boot)",
		),
	)
	return c
}
