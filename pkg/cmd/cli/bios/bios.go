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

package bios

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/stmcginnis/gofish/redfish"
)

// NewCommand creates the `chassis` subcommand.
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:                   "bios",
		DisableFlagsInUseLine: true,
		Short:                 "BIOS control.",
		Long:                  `Interact with a host's BIOS settings.`,
		Hidden:                false,
	}
	c.AddCommand(
		NewGetCommand(),
		NewSetCommand(),
	)
	return c
}

func makeAttributes(args []string) redfish.BiosAttributes {
	attributes := redfish.BiosAttributes{}
	for _, attribute := range args {
		if key, value, ok := strings.Cut(attribute, "="); ok {
			attributes[key] = value
		}
	}
	return attributes
}
