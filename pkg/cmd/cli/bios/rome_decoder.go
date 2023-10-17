/*
 *
 *  MIT License
 *
 *  (C) Copyright 2023 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package bios

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
)

//go:embed rome/*.json
var romeFS embed.FS

type RomeLibrary struct {
	Attributes map[string]Attribute
}

type Attribute struct {
	AttributeName string  `json:"AttributeName"`
	DefaultValue  string  `json:"DefaultValue"`
	DisplayName   string  `json:"DisplayName"`
	HelpText      string  `json:"HelpText"`
	ReadOnly      bool    `json:"ReadOnly"`
	Type          string  `json:"Type"`
	Value         []Value `json:"Value"`
}

type Value struct {
	ValueDisplayName string `json:"ValueDisplayName"`
	ValueName        string `json:"ValueName"`
}

// NewEmbeddedLibrary embeds JSON files from: sh control.Rome.BiosParameters.sh <BMC> renew_json
// new files can be added to rome/ and will be added to the library of available attributes that can be decoded
func NewEmbeddedRomeLibrary(customDir string) (*RomeLibrary, error) {
	library := &RomeLibrary{
		Attributes: map[string]Attribute{},
	}

	// read the embedded directory
	basePath := "rome"
	builtin, err := romeFS.ReadDir(basePath)
	if err != nil {
		os.Exit(1)
	}

	// loop through each json file in the dir
	for _, file := range builtin {
		// read the file's contents
		filePath := path.Join(basePath, file.Name())
		data, err := romeFS.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		// create an Attribute and unmarshal the JSON to it
		attribute := Attribute{}
		err = json.Unmarshal(data, &attribute)
		if err != nil {
			return nil, err
		}

		// add the attribute to the library if it does not yet exist
		_, exists := library.Attributes[attribute.AttributeName]
		if !exists {
			library.Attributes[attribute.AttributeName] = attribute
		}
	}

	return library, nil
}

func (l *RomeLibrary) RegisterAttribute(attribute Attribute) error {
	if _, exists := l.Attributes[attribute.AttributeName]; exists {
		return fmt.Errorf("%s already exists", attribute.AttributeName)
	}

	l.Attributes[attribute.AttributeName] = attribute
	return nil
}
