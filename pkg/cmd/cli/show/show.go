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

package show

import (
	"github.com/spf13/cobra"

	"github.com/Cray-HPE/gru/pkg/cmd/cli/chassis/boot"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/proc"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/system"
)

// NewCommand creates the `show` subcommand.
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "show",
		Short: "Curated server information",
		Long:  `Print pre-defined classes of information from one or more BMCs`,
	}
	c.AddCommand(
		boot.NewShowCommand(),
		proc.NewShowCommand(),
		system.NewShowCommand(),
	)
	return c
}
