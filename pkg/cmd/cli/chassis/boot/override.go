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
	"github.com/Cray-HPE/gru/internal/set"
	"github.com/Cray-HPE/gru/pkg/cmd"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/chassis/power"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stmcginnis/gofish/redfish"
)

// NewBiosOverrideCommand creates the `bios` subcommand for `boot`.
func NewBiosOverrideCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "bios host [...host]",
		Short: "Boot to BIOS",
		Long:  `Override the next boot with the BIOS option`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)

			v := viper.GetViper()
			bindErr := v.BindPFlags(c.Flags())
			cmd.CheckError(bindErr)

			content := set.Async(issueOverride, hosts, redfish.BiosSetupBootSourceOverrideTarget)
			cli.MapPrint(content)
			if v.GetBool("now") {
				content = set.Async(power.Issue, hosts, redfish.ForceRestartResetType)
				cli.MapPrint(content)
			}
		},
	}
	return c
}

// NewPxeOverrideCommand creates the `pxe` subcommand for `boot`.
func NewPxeOverrideCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "pxe host [...host]",
		Short: "Boot with PXE",
		Long:  `Override the next boot with the PXE option`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := set.Async(issueOverride, hosts, redfish.PxeBootSourceOverrideTarget)
			cli.MapPrint(content)
		},
	}
	return c
}

// NewHddOverrideCommand creates the `hdd` subcommand for `boot`.
func NewHddOverrideCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "hdd host [...host]",
		Short: "Boot from the HDD",
		Long:  `Override the next boot with the HDD option`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := set.Async(issueOverride, hosts, redfish.HddBootSourceOverrideTarget)
			cli.MapPrint(content)
		},
	}
	return c
}

// NewUEFIHttpOverrideCommand creates the `http` subcommand for `boot`.
func NewUEFIHttpOverrideCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "http host [...host]",
		Short: "Boot with HTTP",
		Long:  `Override the next boot with the HTTP option`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := set.Async(issueOverride, hosts, redfish.UefiHTTPBootSourceOverrideTarget)
			cli.MapPrint(content)
		},
	}
	return c
}

// NewNoneOverrideCommand creates the `none` subcommand for `boot`.
func NewNoneOverrideCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "none host [...host]",
		Short: "Clear the boot override",
		Long:  `Clears a boot override`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := set.Async(issueOverride, hosts, redfish.NoneBootSourceOverrideTarget)
			cli.MapPrint(content)
		},
	}
	return c
}
