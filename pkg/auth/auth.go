/*
MIT License

(C) Copyright 2022 Hewlett Packard Enterprise Development LP

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
	"github.com/spf13/viper"
	"github.com/stmcginnis/gofish"
	"os"

	"path/filepath"
)

type Configuration struct {
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"IPMI_PASSWORD"`
}

var config Configuration

func LoadConfig(path string) (err error) {
	var cfgFile = viper.GetString("config")
	cfgFile, _ = filepath.Abs(cfgFile)

	if err != nil {
		// FIXME: Don't fail if the env vars exist.
		fmt.Fprintf(os.Stderr, "WARNING: Error reading config file: %v\n", err)
	}

	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	return err
}

// Connection establishes a connection to an endpoint.
func Connection(host string) *gofish.APIClient {
	err := LoadConfig(viper.GetString("config"))
	if err != nil {
		// FIXME: LoadConfig shouldn't fail if env vars were set.
		fmt.Fprintf(os.Stderr, "No credentials given.")
	}
	config := gofish.ClientConfig{
		Endpoint: "https://" + host,
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
		Insecure: viper.GetBool("insecure"),
	}
	c, err := gofish.Connect(config)
	if err != nil {
		panic(err)
	}
	return c
}
