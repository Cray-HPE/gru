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
	"github.com/spf13/cobra"
	"github.com/stmcginnis/gofish/redfish"

	"github.com/Cray-HPE/gru/internal/set"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
)

// NewPowerResetCommand creates the `reset` subcommand for `power`.
func NewPowerResetCommand() *cobra.Command {
	c := &cobra.Command{
		Use: "reset host [...host]",
		Short: "Power reset the target machine(s)",
		Long: `Forcefully restart the target machine(s) without a graceful shutdown`,
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := set.Async(
				Issue,
				hosts,
				redfish.ForceRestartResetType,
			)
			cli.PrettyPrint(content)
		},
		Hidden: false,
	}
	return c
}
