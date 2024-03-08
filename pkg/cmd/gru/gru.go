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

package gru

import (
	"fmt"
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/bios"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/chassis"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/show"
	"github.com/Cray-HPE/gru/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand creates the main command for `gru`.
func NewCommand(name string) *cobra.Command {
	c := &cobra.Command{
		Use:              name,
		TraverseChildren: true,
		Short:            fmt.Sprintf("Go Redfish Utility (%s)", name),
		Long: fmt.Sprintf(
			`
%[1]s is a tool for interacting with Redfish devices. %[1]s provides a
simple interface for interrogating BMC/RMMC Redfish endpoints; reading
and displaying information by pre-defined aggregates, by key-name, and
setting information.

Requires a configuration file or environment variables to be set for 
authenticating with the target RedFish endpoints.

Set USERNAME and PASSWORD (or IPMI_PASSWORD) with the target's credentials, or
provide these via a YAML file. Optionally, if the target hosts have different credentials
the YAML file may provide these per host.
`, name,
		),
		Version: version.Version(),
		PersistentPreRun: func(c *cobra.Command, args []string) {
			v := viper.GetViper()
			bindErr := v.BindPFlags(c.Flags())
			cmd.CheckError(bindErr)
			cfg := v.GetString("config")
			auth.LoadConfig(cfg)
		},
	}
	c.PersistentFlags().StringP(
		"config",
		"c",
		fmt.Sprintf("./%s.yml", name),
		fmt.Sprintln(
			"Configuration file containing BMC credentials, necessary if USERNAME and PASSWORD are not defined in the environment",
			name,
		),
	)
	c.PersistentFlags().Bool(
		"insecure",
		false,
		"Ignore untrusted or insecure certificates",
	)
	c.PersistentFlags().BoolP(
		"json",
		"j",
		false,
		"Output in JSON",
	)
	c.PersistentFlags().BoolP(
		"yaml",
		"y",
		false,
		"Output in YAML",
	)
	c.AddCommand(
		bios.NewCommand(),
		chassis.NewCommand(),
		show.NewCommand(),
	)

	return c
}
