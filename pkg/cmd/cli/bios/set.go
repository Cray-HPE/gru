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

	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/Cray-HPE/gru/pkg/set"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stmcginnis/gofish/redfish"
)

// NewSetCommand creates the `bios` subcommand for `set`.
func NewSetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "bios",
		Short: "Sets BIOS attributes.",
		Args:  cobra.MinimumNArgs(1),
		Long:  `Sets BIOS attributes if the attribute is found and the value is valid.`,
		Run: func(c *cobra.Command, args []string) {
			fromfile, _ = c.PersistentFlags().GetString("from-file")
			defaults, _ = c.PersistentFlags().GetBool("defaults")
			hosts := cli.ParseHosts(args)
			virt = viper.GetBool("virt")
			a := viper.GetStringSlice("attributes")
			attributes := makeAttributes(a)
			content := set.AsyncMap(setBIOSSettings, hosts, attributes)
			cli.MapPrint(content)
		},
		Hidden: false,
	}
	c.PersistentFlags().StringSlice(
		"attributes",
		[]string{},
		"Comma delimited list of attributes and values to set them to: attribute=value[,keyN=valueN].",
	)
	c.PersistentFlags().String(
		"from-file",
		"",
		"Path to a key/value YML file with bios attributes and their desired value(s)",
	)
	c.PersistentFlags().Bool(
		"virt",
		false,
		"Enable virtualization using pre-determined, per-vendor settings",
	)
	c.PersistentFlags().Bool(
		"defaults",
		false,
		"Reset all BIOS attributes to vendor defaults",
	)
	c.MarkFlagsMutuallyExclusive("attributes", "from-file", "virt", "defaults")
	return c
}

// setBIOSSettings sets bios settings on the host
func setBIOSSettings(host string, requestedAttributes map[string]interface{}) interface{} {
	var pendingAttrs = map[string]interface{}{}

	systems, bios, err := getBiosAttributes(host)
	if err != nil {
		pendingAttrs = map[string]interface{}{
			"Error": fmt.Sprintf("%v", err),
		}
		return pendingAttrs
	}

	// Restore default BIOS values
	if defaults {
		err = bios.ResetBios()
		if err != nil {
			pendingAttrs = map[string]interface{}{
				"Error": fmt.Sprintf("BIOS reset failure: %v", err),
			}
			return pendingAttrs
		}
		return pendingAttrs
	}

	var attributes = redfish.BiosAttributes{}

	// if virtualization is requested to be enabled
	if virt {
		// use vendor-specific settings, pre-determined and known to work
		attributes = virtSettings(virt, systems[0].Manufacturer)
	}

	if len(requestedAttributes) == 0 {
		// if no attributes are set, check the file
		if fromfile != "" {
			settings, err := unmarshalBiosKeyValFile(fromfile)
			if err != nil {
				pendingAttrs = map[string]interface{}{
					"Error": fmt.Sprintf("%v", err),
				}
				return pendingAttrs
			}
			// loop through the file's key/values and add them to the map
			// avoid adding duplicate keys
			for k, v := range settings {
				_, exists := requestedAttributes[k]
				if !exists {
					attributes[k] = v
				}
			}
		}
	} else {
		// create a map of all requested attributes
		for k, v := range requestedAttributes {
			attributes[k] = v
		}
	}

	// update the bios settings via redfish
	err = bios.UpdateBiosAttributes(attributes)
	if err != nil {
		pendingAttrs = map[string]interface{}{
			"Error": fmt.Sprintf("could not update BIOS attributes: %v", err),
		}
		return pendingAttrs
	}

	// get the bios object again to check for the pending changes
	_, biosPost, err := getBiosAttributes(host)
	if err != nil {
		pendingAttrs = map[string]interface{}{
			"Error": fmt.Sprintf("%v", err),
		}
		return pendingAttrs
	}

	// get the pending changes
	pendingAttrs = getPendingAttributes(biosPost)

	// return the entire Pending key
	return pendingAttrs
}
