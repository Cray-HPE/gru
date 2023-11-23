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

package boot

import (
	"github.com/Cray-HPE/gru/pkg/auth"
	"github.com/Cray-HPE/gru/pkg/cmd/cli"
	"github.com/spf13/viper"
	"github.com/stmcginnis/gofish/redfish"
)

// issueOverride issues a boot override action against a host.
func issueOverride(host string, override interface{}) interface{} {
	o := cli.Override{}
	v := viper.GetViper()

	c, err := auth.Connection(host)
	if err != nil {
		o.Error = err
		return o
	}

	defer c.Logout()

	service := c.Service

	systems, err := service.Systems()
	if err != nil {
		o.Error = err
		return o
	}

	boot := redfish.Boot{
		BootSourceOverrideTarget: override.(redfish.BootSourceOverrideTarget),
		BootSourceOverrideMode:   redfish.UEFIBootSourceOverrideMode,
	}

	if v.GetBool("persist") {
		boot.BootSourceOverrideEnabled = redfish.ContinuousBootSourceOverrideEnabled
	} else {
		boot.BootSourceOverrideEnabled = redfish.OnceBootSourceOverrideEnabled
	}

	err = systems[0].SetBoot(boot)
	o.Target = override.(redfish.BootSourceOverrideTarget)
	if err != nil {
		o.Error = err
	}

	return o
}
