package bios

import "github.com/stmcginnis/gofish/redfish"

var (
	// IntelEnableVirtualization defines intel-specific bios settings for enabling virtualization
	IntelEnableVirtualization = redfish.SettingsAttributes{
		"VTdSupport":         1,
		"SRIOVEnable":        1,
		"ProcessorX2apic":    1,
		"ProcessorVmxEnable": 1,
	}
	// IntelDisableVirtualization defines intel-specific bios settings for disabling virtualization
	IntelDisableVirtualization = redfish.SettingsAttributes{
		"VTdSupport":         0,
		"SRIOVEnable":        0,
		"ProcessorX2apic":    0,
		"ProcessorVmxEnable": 0,
	}
	// HpeEnableVirtualization hpe-specific bios settings for enabling virtualization
	HpeEnableVirtualization = redfish.SettingsAttributes{
		"ProcAmdVirtualization": "Enabled",
		"ProcAmdIOMMU":          "Enabled",
		"Sriov":                 "Enabled",
		"ProcX2Apic":            "Auto",
	}
	// HpeDisableVirtualization hpe-specific bios settings for disabling virtualization
	HpeDisableVirtualization = redfish.SettingsAttributes{
		"ProcAmdVirtualization": "Disabled",
		"ProcAmdIOMMU":          "Disabled",
		"Sriov":                 "Disabled",
	}
	// GigabyteEnableVirtualization defines gigabyte-specific bios settings for enabling virtualization
	GigabyteEnableVirtualization = redfish.SettingsAttributes{
		"Rome0162": "Enabled", // IOMMU - Enable/Disable IOMMU
		"Rome0565": "Enabled", // SVM Mode - Enable/disable CPU Virtualization
		"PCIS007":  "Enabled", // SR-IOV Support - If system has SR-IOV capable PCIe Devices, this option Enables or Disables Single Root IO Virtualization Support.
		"Rome0059": "Auto",    // SMT Control - Can be used to disable symmetric multithreading. To re-enable SMT, a POWER CYCLE is needed after selecting the 'Auto' option. WARNING - S3 is NOT SUPPORTED on systems where SMT is disabled.
		"Rome0039": "Auto",    // Local APIC Mode - Select local APIC mode: Compatibility, xAPIC or x2APIC
	}
	// GigabyteDisableVirtualization defined gigabyte-specific bios settings for disabling virtualization
	GigabyteDisableVirtualization = redfish.SettingsAttributes{}
)
