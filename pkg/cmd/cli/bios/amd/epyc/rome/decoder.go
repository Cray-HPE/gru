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

package rome

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

//go:embed *.json

var fs embed.FS

// Map is a pointer to a current library of decoded attributes.
var Map *Library

// ProcessorToken is the processor model token.
var ProcessorToken = "^AMD EPYC"

// DecoderMap provides a mapping of attributes for the decoder.
type DecoderMap struct {
	Map *Library
}

// Library is a map of rome bios attributes
type Library struct {
	Attributes map[string]Attribute
}

// Attribute is a single bios attribute
type Attribute struct {
	AttributeName string      `json:"AttributeName"`
	DefaultValue  interface{} `json:"DefaultValue"` // can be int or string, maybe bool
	DisplayName   string      `json:"DisplayName"`
	HelpText      string      `json:"HelpText"`
	ReadOnly      bool        `json:"ReadOnly"`
	Type          string      `json:"Type"`
	Value         []Value     `json:"Value"`
}

// Value is the display name and a name
type Value struct {
	ValueDisplayName string `json:"ValueDisplayName"`
	ValueName        string `json:"ValueName"`
}

// newEmbeddedLibrary embeds JSON files from: sh control.Rome.BiosParameters.sh <BMC> renew_json
// new files can be added to rome/ and will be added to the library of available attributes that can be decoded
func newEmbeddedLibrary(customDir string) (*Library, error) {
	var basePath string
	library := &Library{
		Attributes: map[string]Attribute{},
	}

	if customDir != "" {
		basePath = customDir
	} else {
		basePath = "."
	}
	builtin, err := fs.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range builtin {
		filePath := path.Join(basePath, file.Name())
		data, err := fs.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		attribute := Attribute{}

		err = json.Unmarshal(data, &attribute)
		if err != nil {
			return nil, errors.Join(
				err,
				fmt.Errorf("%+v", string(data)),
			)
		}

		_, exists := library.Attributes[attribute.AttributeName]
		if !exists {
			library.Attributes[attribute.AttributeName] = attribute
		}
	}

	return library, nil
}

// RegisterAttribute adds an attribute to the library
func (l *Library) RegisterAttribute(attribute Attribute) error {
	if _, exists := l.Attributes[attribute.AttributeName]; exists {
		return fmt.Errorf("%s already exists", attribute.AttributeName)
	}

	l.Attributes[attribute.AttributeName] = attribute
	return nil
}

// Decode accepts a key and changes it to a friendly name if it exists and json is not requested
func (d DecoderMap) Decode(key string) string {
	v := viper.GetViper()

	if romeAttr, exists := d.Map.Attributes[key]; exists {
		if v.GetBool("json") {
			key = romeAttr.AttributeName
		} else {
			key = fmt.Sprintf("%s (%s)", romeAttr.AttributeName, strings.TrimLeft(romeAttr.DisplayName, " "))
		}
	}
	return key
}

func init() {
	var err error
	Map, err = newEmbeddedLibrary("")
	if err != nil {
		fmt.Printf("failed to decode rome attributes:\n%v\n", err)
		os.Exit(1)
	}
}
