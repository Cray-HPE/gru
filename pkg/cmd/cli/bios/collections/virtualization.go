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

package collections

import (
	"fmt"
	"github.com/stmcginnis/gofish/redfish"
	"strings"
)

// Virtualization denotes whether to use this collection.
var Virtualization bool

var (
	// IntelEnableVirtualization defines Intel Corporation BIOS settings for enabling virtualization with Intel processors.
	IntelEnableVirtualization = redfish.SettingsAttributes{
		"VTdSupport":         1,
		"SRIOVEnable":        1, // Does not include PCIe NICs, only affects NICs on the physical motherboard.
		"ProcessorX2apic":    1,
		"ProcessorVmxEnable": 1,
	}
	// IntelDisableVirtualization defines Intel Corporation BIOS settings for disabling virtualization with Intel processors.
	IntelDisableVirtualization = redfish.SettingsAttributes{
		"VTdSupport":         0,
		"SRIOVEnable":        0, // Does not include PCIe NICs, only affects NICs on the physical motherboard.
		"ProcessorX2apic":    0,
		"ProcessorVmxEnable": 0,
	}
	// HpeEnableVirtualization Hewlett-Packard Enterprise BIOS settings for enabling virtualization
	HpeEnableVirtualization = redfish.SettingsAttributes{
		"ProcAmdVirtualization": "Enabled",
		"ProcAmdIOMMU":          "Enabled",
		"Sriov":                 "Enabled", // Does not include PCIe NICs, only affects NICs on the physical motherboard.
		"ProcX2Apic":            "Auto",
	}
	// HpeDisableVirtualization Hewlett-Packard Enterprise BIOS settings for disabling virtualization with AMD processors.
	HpeDisableVirtualization = redfish.SettingsAttributes{
		"ProcAmdVirtualization": "Disabled",
		"ProcAmdIOMMU":          "Disabled",
		"Sriov":                 "Disabled", // Does not include PCIe NICs, only affects NICs on the physical motherboard.
		"ProcX2Apic":            "Disabled",
	}
	// GigabyteEnableVirtualization defines Gigabyte Technology BIOS settings for enabling virtualization with AMD processors.
	GigabyteEnableVirtualization = redfish.SettingsAttributes{
		"Rome0162": "Enabled", // IOMMU - Enable/Disable IOMMU
		"Rome0565": "Enabled", // SVM Mode - Enable/disable CPU Virtualization
		"PCIS007":  "Enabled", // SR-IOV Support - If system has SR-IOV capable PCIe Devices, this option Enables or Disables Single Root IO Virtualization Support for onboard NICs! (Does not include PCIe NICs)
		"Rome0059": "Auto",    // SMT Control - Can be used to disable symmetric multithreading. To re-enable SMT, a POWER CYCLE is needed after selecting the 'Auto' option. WARNING - S3 is NOT SUPPORTED on systems where SMT is disabled.
		"Rome0039": "Auto",    // Local APIC Mode - Select local APIC mode: Compatibility, xAPIC or x2APIC
	}
	// GigabyteDisableVirtualization defined Gigabyte Technology bios settings for disabling virtualization with AMD processors.
	GigabyteDisableVirtualization = redfish.SettingsAttributes{
		"Rome0162": "Disabled",
		"Rome0565": "Disabled",
		"PCIS007":  "Disabled",
		"Rome0059": "Disabled",
		"Rome0039": "Disabled",
	}
)

// VirtualizationAttributes is a shortcut for enabling/disabling virtualization as it often
// requires more than one setting to be enabled and varies between vendors
// it takes a bool to determine to enable/disable and the manufacturer name
// to properly determine the settings (disable is currently not supported)
func VirtualizationAttributes(enable bool, manufacturer string) (settings redfish.SettingsAttributes, err error) {
	settings = redfish.SettingsAttributes{}
	normalizedManufacturer := strings.ToUpper(manufacturer)

	switch normalizedManufacturer {
	case "INTEL CORPORATION":
		if enable {
			settings = IntelEnableVirtualization
		} else {
			settings = IntelDisableVirtualization
		}
	case "GIGABYTE", "CRAY INC.":
		if enable {
			settings = GigabyteEnableVirtualization
		} else {
			settings = GigabyteDisableVirtualization
		}
	case "HPE":
		if enable {
			settings = HpeEnableVirtualization
		} else {
			settings = HpeDisableVirtualization
		}
	default:
		err = fmt.Errorf(
			"unable to determine manufaturer for attribute collection; manufacturer detected: %s",
			normalizedManufacturer,
		)
	}
	return settings, err
}
