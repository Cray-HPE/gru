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
	"strings"

	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/stmcginnis/gofish/redfish"
	"gopkg.in/yaml.v3"
)

var (
	fromfile string
	virt     = false
	pending  = false
	defaults = false
)

// getBiosAttributes returns the computer systems, the bios (for system 0), and an error/nil from an endpoint
func getBiosAttributes(host string) (systems []*redfish.ComputerSystem, bios *redfish.Bios, err error) {
	// set up the client
	c, err := auth.Connection(host)
	if err != nil {
		return systems, bios, err
	}
	defer c.Logout()

	// get the sytems
	service := c.Service
	systems, err = service.Systems()
	if err != nil {
		return systems, bios, err
	}

	// get the bios for the first system
	// FIXME: return all?
	bios, err = systems[0].Bios()
	if err != nil {
		return systems, bios, err
	}

	return systems, bios, nil
}

// makeAttributes makes a map[string]interface out of a comma separate slice of \
// key=values, which are split on on '='.  the resultant value is a string,
// which does not work for all attributes
func makeAttributes(args []string) map[string]interface{} {
	attributes := make(map[string]interface{}, 0)
	var a interface{}
	// convert string slice to map breaking on '='
	for _, attribute := range args {
		if key, value, ok := strings.Cut(attribute, "="); ok {
			attributes[key] = value
		}
	}

	// marshal and unmarshal to get types easier
	b, err := yaml.Marshal(attributes)
	if err != nil {
		attributes["Error"] = fmt.Sprintf("%v", err)
	}
	err = yaml.Unmarshal(b, &a)
	if err != nil {
		attributes["Error"] = fmt.Sprintf("%v", err)
	}

	return attributes
}

// unmarshalBiosKeyValFile unmarshals a yaml file into key/value pairs as a map[string]interface{}
func unmarshalBiosKeyValFile(file string) (settings map[string]interface{}, err error) {
	settings = make(map[string]interface{}, 0)

	biosKv, err := os.ReadFile(file)
	if err != nil {
		return settings, err
	}
	err = yaml.Unmarshal(biosKv, settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}

// virtSettings is a shortcut for enabling/disabling virtualization as it often
// requires more than one setting to be enabled and varies between vendors
// it takes a bool to determine to enable/disable and the manufacturer name
// to properly determine the settings (disable is currently not supported)
func virtSettings(enable bool, manufacturer string) (settings redfish.SettingsAttributes) {
	settings = redfish.SettingsAttributes{}

	switch manufacturer {
	case "Intel Corporation":
		if enable {
			settings = IntelEnableVirtualization
		} else {
			return nil
		}
	case "GIGABYTE":
		if enable {
			settings = GigabyteEnableVirtualization
		} else {
			return nil
		}
	case "HPE":
		if enable {
			settings = HpeEnableVirtualization
		} else {
			return nil
		}
	default:
		fmt.Printf("unable to determine manufaturer for virtualization shortcut")
		return settings
	}

	return settings
}
