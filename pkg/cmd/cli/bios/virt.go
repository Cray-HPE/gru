package bios

import "github.com/stmcginnis/gofish/redfish"

var (
	// IntelEnableVirtualization defines intel-specific bios settings for enabling virtualization
	IntelEnableVirtualization = redfish.BiosAttributes{
		"VTdSupport":                     1,
		"SRIOVEnable":                    1,
		"BootMode":                       1,
		"ProcessorX2apic":                1,
		"SvrMngmntAcpiIpmi":              1,
		"ProcessorVmxEnable":             1,
		"ProcessorHyperThreadingDisable": 0,
	}
	// IntelDisableVirtualization defines intel-specific bios settings for disabling virtualization
	IntelDisableVirtualization = redfish.BiosAttributes{
		"VTdSupport":                     0,
		"SRIOVEnable":                    0,
		"BootMode":                       0,
		"ProcessorX2apic":                0,
		"SvrMngmntAcpiIpmi":              0,
		"ProcessorVmxEnable":             0,
		"ProcessorHyperThreadingDisable": 1,
	}
	// HpeEnableVirtualization hpe-specific bios settings for enabling virtualization
	HpeEnableVirtualization = redfish.BiosAttributes{
		"AutoPowerOn":           "AlwaysPowerOff",
		"ProcAmdVirtualization": "Enabled",
		"ProcAmdIOMMU":          "Enabled",
		"Sriov":                 "Enabled",
		"BootMode":              "Uefi",
		"ProcX2Apic":            "Auto",
	}
	// HpeDisableVirtualization hpe-specific bios settings for disabling virtualization
	HpeDisableVirtualization = redfish.BiosAttributes{
		// "AutoPowerOn":           "AlwaysPowerOff",
		"ProcAmdVirtualization": "Disabled",
		"ProcAmdIOMMU":          "Disabled",
		"Sriov":                 "Disabled",
		// "BootMode":              "Uefi",
		// "ProcX2Apic":            "Auto",
	}
	// GigabyteEnableVirtualization defines gigabyte-specific bios settings for enabling virtualization
	GigabyteEnableVirtualization = redfish.BiosAttributes{}
	// GigabyteDisableVirtualization defined gigabyte-specific bios settings for disabling virtualization
	GigabyteDisableVirtualization = redfish.BiosAttributes{}
	// CrayEnableVirtualization defines cray-specific bios settings or enabling virtualization
	CrayEnableVirtualization = redfish.BiosAttributes{}
	// CrayDisableVirtualization defines cray-specific bios settings or disabling virtualization
	CrayDisableVirtualization = redfish.BiosAttributes{}
)
