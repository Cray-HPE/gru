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

	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/Cray-HPE/gru/pkg/query"
	"github.com/spf13/cobra"
	"github.com/stmcginnis/gofish/redfish"
)

// StateChange represents a change in bios settings.
type StateChange struct {
	PreviousPowerState *redfish.Bios           `json:"Bios,omitempty"`
	ResetType          *redfish.BiosAttributes `json:"BiosAttributes,omitempty"`
	Error              error                   `json:"error,omitempty"`
}

// State represents a single bios state.
type State struct {
	Attribute *redfish.Bios `json:"Bios"`
	Error     error         `json:"error,omitempty"`
}

// NewGetCommand creates a `bios` subcommand for `get`.
func NewGetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "get [key[,keyN]]",
		Short: "Gets BIOS settings by key-name, or get every key.",
		Long:  `Gets BIOS settings.`,
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := query.Async(getBIOSSettings, hosts)
			cli.MapPrint(content)
		},
		Hidden: false,
	}
	c.PersistentFlags().StringSlice(
		"attributes",
		[]string{},
		"Comma delimited list of attributes and values to set them to.",
	)
	return c
}

type RecommendedBiosSetting struct {
	Setting string `json:"setting,omitempty"`
	Enabled string `json:"enabled,omitempty"`
	Error   error  `json:"error,omitempty"`
}

type RecommendedBiosSettings []RecommendedBiosSetting

// getBIOSSettings
func getBIOSSettings(host string) interface{} {
	param := RecommendedBiosSetting{}
	c, err := auth.Connection(host)
	if err != nil {
		param.Error = err
		return param
	}
	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil {
		param.Error = err
		return param
	}
	b, err := systems[0].Bios()
	if err != nil {
		param.Error = err
		return param
	}

	params := RecommendedBiosSettings{}

	for a, _ := range b.Attributes {
		if a == "AutoPowerOn" ||
			a == "VTdSupport" || a == "ProcAmdVirtualization" ||
			a == "SRIOVEnable" || a == "Sriov" ||
			a == "ProcessorHyperThreadingDisable" ||
			a == "SvrMngmntAcpiIpmi" || a == "ProcX2Apic" {
			state := b.Attributes.Bool(a)
			param.Setting = a
			param.Enabled = fmt.Sprintf("%v", state)
			params = append(params, param)
		}
	}

	return *b
}

// // issue issues an action against a host.
// func issue(host string, action interface{}) interface{} {
// 	sc := StateChange{}
// 	c, err := auth.Connection(host)
// 	if err != nil {
// 		sc.Error = err
// 		return sc
// 	}
// 	defer c.Logout()

// 	service := c.Service

// 	systems, err := service.Systems()
// 	if err != nil {
// 		sc.Error = err
// 		return sc
// 	}
// 	sc.PreviousPowerState = systems[0].PowerState
// 	sc.ResetType = action.(redfish.ResetType)
// 	err = systems[0].Reset(sc.ResetType)
// 	if err != nil {
// 		sc.Error = err
// 	}

// 	return sc
// }
