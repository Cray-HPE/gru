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
	"github.com/spf13/cobra"
)

// NewSetCommand creates the `bios` subcommand for `set`.
func NewSetCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "set attribute=value[,keyN=valueN]",
		Short: "Sets BIOS attributes.",
		Long:  `Sets BIOS attributes if the attribute is found and the value is valid.`,
		Run: func(c *cobra.Command, args []string) {
			// hosts := cli.ParseHosts(args)
			// a := viper.GetStringSlice("attributes")
			// attributes := makeAttributes(a)
			// content := set.AsyncMap(setBIOSSettings, hosts, attributes)
			content := getBIOSSettings(args[0])
			fmt.Print(content)
			// cli.MapPrint(content)
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

// setBIOSSettings
func setBIOSSettings(host string, attributes map[string]interface{}) interface{} {
	c, err := auth.Connection(host)
	if err != nil {
		return err
	}
	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil {
		return err
	}
	b, err := systems[0].Bios()
	if err != nil {
		return b
	}
	// err = bios.UpdateBiosAttributes(attributes)
	// if err != nil {
	// 	return bios
	// }

	for b, i := range b.Description {
		fmt.Printf("%v-----------%+v\n", i, b)
	}
	return b
}
