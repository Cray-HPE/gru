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

package auth

import (
	"fmt"

	"github.com/Cray-HPE/gru/pkg/cmd"
	"github.com/spf13/viper"
	"github.com/stmcginnis/gofish"
)

type configuration struct {
	username string
	password string
	hosts    map[string]interface{}
}

var config configuration

// LoadConfig loads the applications configuration file and merges it with the environment.
func LoadConfig(path string) {
	viper.SetDefault("username", "")
	viper.SetDefault("password", "")
	if viper.BindEnv("password", "IPMI_PASSWORD") != nil {
		cmd.CheckError(fmt.Errorf("failed to bind ipmi_password environment variable"))
	}
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Printf("Loading config file %s", path)
		} else {
			// Config file was found but another error was produced
			// TODO: Handle all errors except if the file is missing.
		}
	}

	viper.AutomaticEnv()
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}

// Connection establishes a connection to an endpoint.
func Connection(host string) (*gofish.APIClient, error) {
	if (viper.GetString("username") == "") || (viper.GetString("password") == "") {
		cmd.CheckError(fmt.Errorf("no credentials provided, please provide a config file or environment variables"))
	}
	hosts := viper.GetStringMap("hosts")
	username := viper.GetString("username")
	password := viper.GetString("password")
	if val, ok := hosts[host]; ok {
		hostConfig := val.(map[string]interface{})
		username = fmt.Sprintf("%v", hostConfig["username"])
		password = fmt.Sprintf("%v", hostConfig["password"])
	}
	config := gofish.ClientConfig{
		Endpoint:  "https://" + host,
		Username:  username,
		Password:  password,
		Insecure:  viper.GetBool("insecure"),
		BasicAuth: true,
	}
	c, err := gofish.Connect(config)
	return c, err
}
