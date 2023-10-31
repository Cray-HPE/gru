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

package cli

import (
	"github.com/stmcginnis/gofish/redfish"
)

// StateChange represents a change in power states.
type StateChange struct {
	PreviousPowerState redfish.PowerState `json:"previousPowerState,omitempty"`
	ResetType          redfish.ResetType  `json:"resetType,omitempty"`
	Error              error              `json:"error,omitempty"`
}

// State represents a single power state.
type State struct {
	PowerState redfish.PowerState `json:"powerState"`
	Error      error              `json:"error,omitempty"`
}

// Boot represents boot configuration on the BMC. Only Error is emitted on empty.
type Boot struct {
	Order []BootDescription `json:"order"`
	Next  string            `json:"next"`
	Error error             `json:"error,omitempty"`
}

type BootDescription map[string]string

// Override represents the result of the boot override.
type Override struct {
	Target redfish.BootSourceOverrideTarget `json:"target"`
	Error  error                            `json:"error,omitempty"`
}
