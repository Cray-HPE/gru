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

package proc

import (
	"fmt"
	"github.com/Cray-HPE/gru/internal/query"
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/spf13/cobra"
	"strings"
)

// NewShowCommand creates the `system` subcommand for `show`.
func NewShowCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "proc host [...host]",
		Short: "Processor information",
		Long:  `Show the Server's processors, a full list with their core count, model, architecture, and serial numbers.`,
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := query.Async(getProcessors, hosts)
			cli.PrettyPrint(content)
		},
	}
	return c
}

func getProcessors(host string) interface{} {
	foundProcessors := Processors{}
	c, err := auth.Connection(host)
	if err != nil {
		foundProcessors = append(
			foundProcessors, Processor{
				Error: err,
			},
		)
		return foundProcessors
	}
	defer c.Logout()
	service := c.Service

	managers, err := service.Managers()
	if err != nil || len(managers) < 1 {
		foundProcessors = append(
			foundProcessors, Processor{
				Error: err,
			},
		)
		return foundProcessors
	}

	systems, err := service.Systems()
	if err != nil || len(systems) < 1 {
		foundProcessors = append(
			foundProcessors, Processor{
				Error: err,
			},
		)
		return foundProcessors
	}
	systemProcessors, err := systems[0].Processors()
	if err != nil {
		foundProcessors = append(
			foundProcessors, Processor{
				Error: err,
			},
		)
	} else {
		for i := range systemProcessors {
			processor := Processor{
				Architecture: strings.TrimSpace(fmt.Sprintf("%v", systemProcessors[i].ProcessorArchitecture)),
				TotalCores:   systemProcessors[i].TotalCores,
				Model:        strings.TrimSpace(systemProcessors[i].Model),
				Socket:       strings.TrimSpace(systemProcessors[i].Socket),
				Threads:      systemProcessors[i].TotalThreads,
				VendorID:     strings.TrimSpace(systemProcessors[i].ProcessorID.VendorID),
			}
			foundProcessors = append(foundProcessors, processor)
		}
	}
	return foundProcessors
}
