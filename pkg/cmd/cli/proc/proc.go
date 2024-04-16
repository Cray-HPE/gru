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

package proc

// Processors represents a list of Processor types.
type Processors []Processor

// Processor represents a single, physical processor.
// The processor's model number is stored in either Processor.Model or Processor.VendorID depending
// on the vendor. The user may need to interpret both fields to understand what they have.
type Processor struct {
	Architecture string `json:"architecture" yaml:"architecture"`
	TotalCores   int    `json:"totalCores" yaml:"total_cores"`
	Model        string `json:"model" yaml:"model"`
	Socket       string `json:"socket" yaml:"socket"`
	Threads      int    `json:"threads" yaml:"threads"`
	VendorID     string `json:"vendorID" yaml:"vendor_id"`
	Error        error  `json:"error,omitempty" yaml:"error,omitempty"`
}
