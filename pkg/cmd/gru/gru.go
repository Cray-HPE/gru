/*
MIT License

(C) Copyright 2022 Hewlett Packard Enterprise Development LP

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
	"github.com/Cray-HPE/gru/pkg/cmd/cli/get"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/show"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/version"
	"github.com/spf13/cobra"
)

func NewCommand(name string) *cobra.Command {
	c := &cobra.Command{
		Use:   name,
		Short: "Go Redfish Utility (gru)",
		Long: `
gru is a tool for reading and writing to BMCs via Redfish. It provides a
simple interface for interrogating BMCs; reading and displaying information by pre-defined aggregates,
by key-name, and setting information.
`,
	}
	c.PersistentFlags().StringP("config", "c", "gru.yaml", "Configuration file containing BMC credentials, necessary if `USERNAME` and `IPMI_PASSWORD` is not defined in the environment.")
	c.AddCommand(
		get.NewCommand(),
		show.NewCommand(),
		version.NewCommand(name),
	)
	return c
}
