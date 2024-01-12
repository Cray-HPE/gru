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
	"encoding/json"
	"fmt"
	"github.com/Cray-HPE/gru/internal/query"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/bios/collections"
	"github.com/spf13/cobra"
	"github.com/stmcginnis/gofish/redfish"
	"log"
	"regexp"

	"github.com/spf13/viper"
	"strings"
)

// NewBiosGetCommand creates a `get` subcommand for `bios`.
func NewBiosGetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "get",
		Short: "Gets BIOS attributes by key-name, or get all attributes",
		Long:  `Gets BIOS attributes`,
		Run: func(c *cobra.Command, args []string) {
			hosts := cli.ParseHosts(args)
			content := query.Async(getBiosAttributes, hosts)
			cli.MapPrint(content)
		},
		Hidden: false,
	}

	c.PersistentFlags().BoolP(
		"pending",
		"p",
		false,
		"Get pending BIOS attribute changes",
	)
	return c
}

// getBiosAttributes gets the requested attribute names or gets all attributes
func getBiosAttributes(host string) interface{} {
	v := viper.GetViper()
	var biosDecoder Decoder
	var requestedAttributes []string
	attributes := Attributes{}

	if v.GetBool("pending") {
		pendingAttributes := getPendingBiosAttributes(host)
		attributes.Pending = pendingAttributes.Pending
		attributes.Error = pendingAttributes.Error
		return attributes
	}

	systems, bios, err := getSystemBios(host)
	if err != nil {
		return err
	}

	for decoder := range AttributeDecoderMaps {
		regex, err := regexp.Compile(AttributeDecoderMaps[decoder].Token)
		if err != nil {
			fmt.Printf("Failed to create decoder regex for: %s", AttributeDecoderMaps[decoder].Token)
			continue
		}
		if regex.MatchString(systems[0].ProcessorSummary.Model) {
			biosDecoder = AttributeDecoderMaps[decoder]
			break
		}
	}

	// TODO: retry/timeout N before failing
	if len(bios.Attributes) == 0 {
		attributes.Error = fmt.Errorf("node may be off, or in a broken state, or unrecognizeable by gru")
		return attributes
	}

	// if the values are coming from a file
	fromFile := v.GetString("from-file")
	if fromFile != "" {
		attrsFromFile, err := unmarshalBiosKeyValFile(fromFile)
		if err != nil {
			log.Fatal(err)
		}
		for k := range attrsFromFile {
			requestedAttributes = append(requestedAttributes, k)
		}
	} else {
		// get the keys from the attributes flags
		requestedAttributes = viper.GetStringSlice("attributes")
	}

	if v.GetBool("virtualization") {
		virtualizationAttributes, err := collections.VirtualizationAttributes(true, systems[0].Manufacturer)
		if err != nil {
			attributes.Error = err
			return attributes
		}

		for key := range virtualizationAttributes {
			requestedAttributes = append(requestedAttributes, key)
		}
	}

	attributes.Attributes = redfish.SettingsAttributes{}

	if len(requestedAttributes) != 0 {
		// check each requested attribute
		for _, attribute := range requestedAttributes {
			if biosDecoder != nil {
				attribute = biosDecoder.Decode(attribute)
			}
			if v, exists := bios.Attributes[attribute]; exists {
				attributes.Attributes[attribute] = v
			} else {
				attributes.Attributes[attribute] = nil
			}
		}
		if len(attributes.Attributes) == 0 {
			attributes.Error = fmt.Errorf("no matching keys found in: %v", requestedAttributes)
		}
	} else {

		for k, v := range bios.Attributes {
			if biosDecoder != nil {
				k = biosDecoder.Decode(k)
				attributes.Attributes[k] = v
			} else {
				attributes.Attributes[k] = v
			}
		}
	}
	return attributes
}

// getPendingBiosAttributes gets the staged bios attributes from Bios/Settings
func getPendingBiosAttributes(host string) Attributes {
	attributes := Attributes{}

	_, bios, err := getSystemBios(host)
	if err != nil {
		attributes.Error = err
		return attributes
	}

	/*
		Redfish will stage the changes in "Bios/Settings"
		Check if it exists before declaring success
		Some combos of redfish/bios versions do not actually have this endpoint
		The library should actually check for this, but this works for now
	*/
	staging := fmt.Sprintf("%s/%s", strings.TrimRight(bios.ODataID, "/"), "Settings")
	resp, err := bios.Client.Get(staging)
	if err != nil {
		attributes.Error = err
		return attributes
	}

	/*
		make a simple map for checking the existence of the "Attributes" key
		if it does not exist, ``bios.UpdateBiosAttributes`` still returns 200
		even though no changes can actually take place
	*/
	staged := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&staged)
	if err != nil {
		attributes.Error = err
		return attributes
	}

	_, exists := staged["Attributes"]
	if staged["Attributes"] == nil || !exists {
		attributes.Error = fmt.Errorf("\"Attributes\" does not exist or is null, the BIOS/firmware may need to updated for proper Attributes support")
		return attributes
	}

	modified := make(map[string]interface{})
	for k, v := range staged["Attributes"].(map[string]interface{}) {
		if v != bios.Attributes[k] {
			modified[k] = v
		}
	}

	// TODO: might want to add the ApplyTimes so the user knows when the change could take effect
	attributes.Pending = modified
	return attributes
}
