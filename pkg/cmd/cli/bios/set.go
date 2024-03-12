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

package bios

import (
	"fmt"
	"github.com/Cray-HPE/gru/internal/set"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/bios/collections"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stmcginnis/gofish/redfish"
	"os"
)

// ClearCmos determines whether to clear the CMOS values.
var ClearCmos bool

// NewBiosSetCommand creates the `set` subcommand for `bios`.
func NewBiosSetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "set host [...host]",
		Short: "Sets BIOS attributes",
		Long:  `Sets BIOS attributes if the attribute is found and the value is valid.`,
		Run: func(c *cobra.Command, args []string) {
			if len(Attributes) == 0 && FromFile == "" && !collections.Virtualization && !ClearCmos {
				_, err := fmt.Fprintln(
					os.Stderr,
					fmt.Errorf("an error occurred: at least one of the flags in the group [attributes from-file virtualization clear-cmos] is required"),
				)
				err = c.Help()
				if err != nil {
					return
				}
				os.Exit(1)
			}
			if (len(Attributes) == 0) == (FromFile == "") == collections.Virtualization == ClearCmos {
				_, err := fmt.Fprintln(
					os.Stderr,
					fmt.Errorf("an error occurred: only one of the flags in the group [attributes from-file virtualization clear-cmos] can be specified at a time"),
				)
				err = c.Help()
				if err != nil {
					return
				}
				os.Exit(1)
			}

			v := viper.GetViper()
			hosts := cli.ParseHosts(args)
			a := viper.GetStringSlice("attributes")
			attributes := makeAttributes(a)
			var content map[string]interface{}

			if v.GetBool("clear-cmos") {
				content = set.AsyncCall(resetBios, hosts)
			} else {
				content = set.AsyncMap(setBios, hosts, attributes.Attributes)
			}

			cli.PrettyPrint(content)
		},
		Hidden: false,
	}

	c.PersistentFlags().BoolVar(
		&ClearCmos,
		"clear-cmos",
		false,
		"Clear CMOS; set all BIOS attributes to their defaults.",
	)

	return c
}

func setBios(host string, requestedAttributes map[string]interface{}) interface{} {
	attributes := Settings{}
	v := viper.GetViper()

	systems, bios, err := getSystemBios(host)
	if err != nil || len(systems) < 1 {
		attributes.Error = err
		return attributes
	}

	attributes.Attributes = redfish.SettingsAttributes{}

	if v.GetBool("virtualization") {
		attributes.Attributes, err = collections.VirtualizationAttributes(true, systems[0].Manufacturer)
		if err != nil {
			attributes.Error = err
			return attributes
		}
	}

	if len(requestedAttributes) == 0 {

		fromFile := v.GetString("from-file")

		if fromFile != "" {

			settings, err := unmarshalBiosKeyValFile(fromFile)
			if err != nil {
				attributes.Error = err
				return attributes
			}

			for k, v := range settings {
				_, exists := requestedAttributes[k]
				if !exists {
					attributes.Attributes[k] = v
				}
			}

		}
	} else {

		for k, v := range requestedAttributes {
			attributes.Attributes[k] = v
		}

	}

	err = bios.UpdateBiosAttributes(attributes.Attributes)
	if err != nil {
		attributes.Error = err
		return attributes
	}

	pendingAttributes := getPendingBiosAttributes(host)
	attributes.Pending = pendingAttributes.Pending

	return attributes
}

func resetBios(host string) interface{} {
	attributes := Settings{}

	_, bios, err := getSystemBios(host)
	if err != nil {
		attributes.Error = err
		return attributes
	}

	err = bios.ResetBios()
	if err != nil {
		attributes.Error = err
		return attributes
	}

	return attributes
}
