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
	"os"
	"strings"

	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/Cray-HPE/gru/pkg/query"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stmcginnis/gofish/redfish"
)

// NewGetCommand creates a `bios` subcommand for `get`.
func NewGetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "bios",
		Short: "Gets BIOS settings by key-name, or get every key.",
		Long:  `Gets BIOS settings.`,
		// Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			fromfile, _ = c.PersistentFlags().GetString("from-file")
			pending, _ = c.PersistentFlags().GetBool("pending")
			virt = viper.GetBool("virt")
			hosts := cli.ParseHosts(args)
			content := query.Async(getBIOSSettings, hosts)
			cli.MapPrint(content)
		},
		Hidden: false,
	}
	c.PersistentFlags().StringSlice(
		"attributes",
		[]string{},
		"Comma delimited list of attributes and values: [key[,keyN]]",
	)
	c.PersistentFlags().String(
		"from-file",
		"",
		"Path to a key/value YML file with bios attributes (value(s) for key(s) will be ignored)",
	)
	c.PersistentFlags().Bool(
		"pending",
		false,
		"Get pending BIOS attribute changes",
	)
	c.PersistentFlags().Bool(
		"virt",
		false,
		"Shortcut to get all pre-determined, per-vendor settings for virtualization",
	)
	return c
}

// getBIOSSettings gets the requested attribute names or gets all attributes
func getBIOSSettings(host string, requestedAttributes ...string) interface{} {
	systems, bios, err := getBiosAttributes(host)
	if err != nil {
		return err
	}

	// if the pending attributes are requested, get those and return
	if pending {
		pendingAttrs := getPendingAttributes(bios)
		return pendingAttrs
	}

	// some nodes do not report this endpoint if the node is off
	// TODO: retry/timeout N before failing
	if len(bios.Attributes) == 0 {
		return fmt.Errorf("this node may be off or in a broken state (a check is not yet implemented for this)")
	}

	// if the values are coming from a file
	if fromfile != "" {
		// unmarshal
		attrsFromFile, err := unmarshalBiosKeyValFile(fromfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// only the key name is needed for GET requests
		// this allows the same file to be re-used for get and set since the
		// value is just ignored here
		for k := range attrsFromFile {
			requestedAttributes = append(requestedAttributes, k)
		}
	} else {
		// get the keys from the attributes flags
		requestedAttributes = viper.GetStringSlice("attributes")
	}

	attributes := redfish.SettingsAttributes{}
	if len(requestedAttributes) != 0 {
		// check each requested attribute
		for _, attribute := range requestedAttributes {
			// check if it the requested key exists
			v, exists := bios.Attributes[attribute]
			if exists {
				// add it to the map if it does
				attributes[attribute] = v
			} else {
				attributes[attribute] = nil
			}
		}
		// fail if no matching keys were found
		if len(attributes) == 0 {
			return fmt.Errorf("no matching keys found in: %v", requestedAttributes)
		}
	} else {
		// if virtualization is requested to be enabled
		if virt {
			// use vendor-specific settings, pre-determined and known to work
			virtAttrs := virtSettings(virt, systems[0].Manufacturer)
			for k := range virtAttrs {
				romeAttr, exists := romeMap.Attributes[k]
				if exists {
					// convert to a friendly name for non-json
					if viper.GetBool("json") {
						k = romeAttr.AttributeName
					} else {
						k = fmt.Sprintf("%s (%s)", romeAttr.AttributeName, romeAttr.DisplayName)
					}
				}
				attributes[k] = bios.Attributes[k]
			}
		} else {
			// loop through all keys discovered and add them to the returned map
			for k, v := range bios.Attributes {
				attributes[k] = v
			}
		}
	}

	return attributes
}

// getPendingAttributes gets the staged bios attributes from Bios/Settings
func getPendingAttributes(bios *redfish.Bios) map[string]interface{} {
	var pendingAttrs = map[string]interface{}{}

	// Redfish will "stage" the changes in "Bios/Settings"
	// Check if it exists before declaring success
	// Some combos of redfish/bios versions do not actually have this endpoint
	// The library should actually check for this, but this works for now
	staging := fmt.Sprintf("%s/%s", strings.TrimRight(bios.ODataID, "/"), "Settings")
	resp, err := bios.Client.Get(staging)
	if err != nil {
		pendingAttrs = map[string]interface{}{
			"Error": fmt.Sprintf("%v", err),
		}
		return pendingAttrs
	}

	// make a simple map for checking the existence of the "Attributes" key
	// if it does not exist, bios.UpdateBiosAttributes still returns 200
	// even though no changes can actually take place
	staged := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&staged)
	if err != nil {
		pendingAttrs = map[string]interface{}{
			"Error": fmt.Sprintf("%v", err),
		}
		return pendingAttrs
	}

	// if it does not exist, it could indicate the bios or firmware need updates
	_, exists := staged["Attributes"]
	if staged["Attributes"] == nil || !exists {
		err := fmt.Errorf("\"Attributes\" does not exist or is null.  You may need to update the BIOS/firmware")
		pendingAttrs = map[string]interface{}{
			"Error": fmt.Sprintf("%v", err),
		}
		return pendingAttrs
	}

	// get staged modified keys
	modified := make(map[string]interface{})
	for k, v := range staged["Attributes"].(map[string]interface{}) {
		// if the value does not match the value of the current setting, add it to
		// the map to show only the modified ones
		if v != bios.Attributes[k] {
			modified[k] = v
		}
	}
	// return a Pending key, this is not a redfish construct, but one that is nice for humans and robots
	// TODO: might want to add the ApplyTimes so the user knows when the change could take effect
	pendingAttrs = map[string]interface{}{
		"Pending": map[string]interface{}{
			"Attributes": modified,
		},
	}
	return pendingAttrs
}
