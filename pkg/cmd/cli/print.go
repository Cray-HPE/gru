/*

 MIT License

 (C) Copyright 2024 Hewlett Packard Enterprise Development LP

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

package cli

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"reflect"
	"sort"
)

// keyValuePrint is a print helper for formatting key-value pairs.
func keyValuePrint(key string, value any, indent string) {
	fmt.Printf("%s%-60s: %-60v\n", indent, key, value)
}

// keyPrint is a print helper for printing just a key in anticipation of printing a data structure as a value.
func keyPrint(key string, indent string) {
	fmt.Printf("%s%s:\n", indent, key)
}

// sliceElementPrint is a print helper for printing values belonging to an array in a marked up format.
func sliceElementPrint(value any, indent string) {
	fmt.Printf("%s%-60v\n", indent, value)
}

// slicePrint is a print helper for printing a slice of interfaces.
func slicePrint(value reflect.Value) {

	for i := 0; i < value.Len(); i++ {

		sv := reflect.ValueOf(value.Index(i).Interface())
		keyPrint(fmt.Sprintf("%d", i), "\t")

		for k := 0; k < sv.NumField(); k++ {

			typeOfS := sv.Type()

			if sv.Field(k).Interface() == nil {

				continue

			}

			kind := reflect.TypeOf(sv.Field(k).Interface()).Kind()

			if kind == reflect.Slice {

				keyPrint(typeOfS.Field(k).Name, "\t")

				// FIXME: This does not increase indent level.
				slicePrint(sv.Field(k))

			} else {
				keyValuePrint(typeOfS.Field(k).Name, sv.Field(k), "\t\t")
			}
		}
	}
}

// structPrint is a print helper for printing a map.
func structPrint(value reflect.Value) {

	typeOfS := value.Type()

	for i := 0; i < value.NumField(); i++ {

		if value.Field(i).Interface() == nil {

			continue

		} else if _, ok := value.Field(i).Interface().(map[string]interface{}); ok {

			keys := value.Field(i).MapKeys()

			if len(keys) != 0 {

				keyPrint(typeOfS.Field(i).Name, "\t")

			} else {

				continue

			}

			sortedKeys := make([]string, 0, len(keys))

			for key := range keys {

				sortedKeys = append(sortedKeys, keys[key].String())

			}

			sort.Strings(sortedKeys)

			for key := range sortedKeys {

				keyValuePrint(sortedKeys[key], value.Field(i).MapIndex(reflect.ValueOf(sortedKeys[key])), "\t\t")

			}

		} else if _, ok := value.Field(i).Interface().([]string); ok {

			keyPrint(typeOfS.Field(i).Name, "\t")

			for _, v := range value.Field(i).Interface().([]string) {

				sliceElementPrint(v, "\t\t")

			}

		} else {
			if value.Field(i).Interface() == nil || value.Field(i).Interface() == "" {

				continue

			}
			if reflect.TypeOf(value.Field(i).Interface()).Kind() == reflect.Slice {

				slicePrint(value.Field(i))

			}

			keyValuePrint(typeOfS.Field(i).Name, value.Field(i).Interface(), "\t")

		}
	}
}

// PrettyPrint prints output in a human-readable manner, unless a specific format is given (e.g. `--json` or `--yaml`).
func PrettyPrint(content map[string]interface{}) {
	if viper.GetBool("json") {

		JSON, err := json.MarshalIndent(content, "", "  ")
		if err != nil {
			panic(fmt.Errorf("could not create valid JSON from %v", content))
		}
		fmt.Printf("%s\n", string(JSON))

	} else if viper.GetBool("yaml") {

		YAML, err := yaml.Marshal(content)
		if err != nil {
			panic(fmt.Errorf("could not create valid YAML from %v", content))
		}
		fmt.Printf("%s\n", string(YAML))

	} else {

		keys := make([]string, 0, len(content))

		for k := range content {

			keys = append(keys, k)

		}

		sort.Strings(keys)

		for _, k := range keys {

			fmt.Printf("%s:\n", k)

			// Warning; the struct fields must be exported!
			s := content[k]
			v := reflect.ValueOf(s)

			if reflect.TypeOf(v.Interface()).Kind() == reflect.Slice {

				slicePrint(v)

			} else if reflect.TypeOf(v.Interface()).Kind() == reflect.Struct {

				structPrint(v)

			} else {

				fmt.Println("Shouldn't be here.")

			}
		}
	}
}
