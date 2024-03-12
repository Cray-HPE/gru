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
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd/cli/bios/collections"
	"github.com/spf13/cobra"
	"github.com/stmcginnis/gofish/redfish"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Settings is a structure for holding current BIOS attributes, pending attributes, and errors.
type Settings struct {
	Attributes map[string]interface{} `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	Pending    map[string]interface{} `json:"pending,omitempty" yaml:"pending,omitempty"`
	Error      error                  `json:"error,omitempty" yaml:"error,omitempty"`
}

// Attributes are an array of attribute names (and optionally values).
var Attributes []string

// FromFile is a path to a file to read attributes from.
var FromFile string

// NewCommand creates the `bios` subcommand.
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:              "bios",
		Short:            "BIOS interaction",
		Long:             `Interact with a host's bios`,
		TraverseChildren: true,
		Hidden:           false,
		Run: func(c *cobra.Command, args []string) {
		},
	}

	c.PersistentFlags().StringArrayVarP(
		&Attributes,
		"attributes",
		"a",
		[]string{},
		"Comma delimited list of attributes and values: [key[,keyN]]",
	)

	c.PersistentFlags().StringVarP(
		&FromFile,
		"from-file",
		"f",
		"",
		"Path to an INI or YAML file with bios attributes (value(s) for key(s) will be ignored)",
	)

	c.PersistentFlags().BoolVarP(
		&collections.Virtualization,
		"virtualization",
		"V",
		false,
		"Shortcut to get all pre-determined, per-vendor settings for virtualization",
	)

	c.AddCommand(
		NewBiosGetCommand(),
		NewBiosSetCommand(),
	)
	return c
}

// makeAttributes makes a “map[string]interface“ out of a comma separate slice of
// key=values, which are split on on '='.  the resultant value is a string,
// which does not work for all attributes.
func makeAttributes(args []string) Settings {
	attributes := Settings{}
	attributes.Attributes = make(map[string]interface{})

	var a interface{}

	for _, attribute := range args {
		if key, value, ok := strings.Cut(attribute, "="); ok {
			attributes.Attributes[key] = value
		}
	}

	b, err := yaml.Marshal(attributes.Attributes)
	if err != nil {
		attributes.Error = err
		return attributes
	}

	err = yaml.Unmarshal(b, &a)
	if err != nil {
		attributes.Error = err
		return attributes
	}

	return attributes
}

// unmarshalBiosKeyValFile unmarshal an INI or YAML file into key/value pairs as a “map[string]interface{}“.
func unmarshalBiosKeyValFile(file string) (settings map[string]interface{}, err error) {
	settings = make(map[string]interface{})

	biosKv, err := os.ReadFile(file)
	if err != nil {
		return settings, err
	}

	fileExtension := filepath.Ext(file)
	re := regexp.MustCompile(`ya?ml`)
	if re.Match([]byte(fileExtension)) {
		err = yaml.Unmarshal(biosKv, settings)
	} else {
		return settings, fmt.Errorf("invalid filetype: %s", fileExtension)
	}

	if err != nil {
		return settings, err
	}
	return settings, nil
}

// getSystemBios returns a slice of “redfish.ComputerSystem“, “redfish.Bios“, and an “error“ (if there was one)
// from an endpoint as-is, we get all systems but only return system[0].Bios, there could be different Bios per system
// someone could also use the wrong system in the returned slice of systems.
// TODO: return a map of systems to bios objects
func getSystemBios(host string) (systems []*redfish.ComputerSystem, bios *redfish.Bios, err error) {

	c, err := auth.Connection(host)
	if err != nil {
		return systems, bios, err
	}
	defer c.Logout()

	// get the systems
	service := c.Service
	systems, err = service.Systems()
	if err != nil || len(systems) < 1 {
		return systems, bios, err
	}

	// TODO from above: create map[string]*redfish.Bios (systems[0].HostName)
	bios, err = systems[0].Bios()
	if err != nil {
		return systems, bios, err
	}

	return systems, bios, nil
}
