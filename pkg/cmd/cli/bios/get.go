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
	"fmt"
	"os"
	"sort"

	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/Cray-HPE/gru/pkg/query"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stmcginnis/gofish/redfish"
)

// NewGetCommand creates a `bios` subcommand for `get`.
func NewGetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "get [key[,keyN]]",
		Short: "Gets BIOS settings by key-name, or get every key.",
		Long:  `Gets BIOS settings.`,
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			err := query.Async(getBIOSSettings, hosts)
			if err != nil {
				os.Exit(1)
			}
			// cli.MapPrint(content)
		},
		Hidden: false, // TODO: Remove or set to false once implemented.
	}
	c.PersistentFlags().StringSlice(
		"attributes",
		[]string{},
		"Comma delimited list of attributes and values to set them to.",
	)
	return c
}

type Bios struct {
	Attributes redfish.BiosAttributes
	// Simulator:
	//
	// BmcAdminMode
	// ProcessorHyperThreadingDisable
	// SRIOVEnable
	// SvrMngmntResumeAcPowerLoss
	// VTdSupport

	// HPE:
	//
	// AutoPowerOn
	// ProcAmdVirtualization:
	// ProcX2Apic

	// Intel:
	//
	// ProcessorX2apic

}

// FIXME: This is a skeleton, and is neither done nor correct. It is a napkin of how this could work.
func getBIOSSettings(host string) interface{} {
	c, err := auth.Connection(host)
	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	allBios := map[*redfish.ComputerSystem]*redfish.BiosAttributes{}
	for _, s := range systems {
		bios, err := s.Bios()
		if err != nil {
			fmt.Println(err)
			return nil
		}

		// Add the key names to a slice
		sorted := []string{}
		for attr, _ := range bios.Attributes {
			sorted = append(sorted, attr)
		}

		// Sort it
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i] < sorted[j]
		})

		a := viper.GetStringSlice("attributes")
		if len(a) != 0 {
			for _, key := range a {

				// FIXME: There is no validation for the key, and no handling if for the key if it doesn't exist.
				// FIXME: Should handle vendor specific items where feasible, translating items such as "VTT" to the vendor's BIOS option.

				state := bios.Attributes.Bool(key)
				fmt.Printf("%v: %v\n", key, state)
			}
		} else {
			// Print all keys and their state
			for _, attr := range sorted {
				state := bios.Attributes.Bool(attr)
				fmt.Printf("%v: %v\n", attr, state)
			}
		}

	}
	return allBios
}
