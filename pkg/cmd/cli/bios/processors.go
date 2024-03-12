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

package bios

/* Maintainer Note:
Every new decoder will need to be imported, and added to `var AttributeDecoderMaps`.
*/
import (
	"github.com/Cray-HPE/gru/pkg/cmd/cli/bios/amd/epyc/rome"
)

// Decoder is an interface for decoding keys (strings) against translated BIOS attributes.
type Decoder interface {
	Decode(key string) string
}

// DecoderMap maps a token to a decoder.
type DecoderMap struct {
	Token   string
	Decoder Decoder
}

// Decode decodes a given key against a DecoderMap of tokens and returns the decoded string.
func (d *DecoderMap) Decode(key string) string {
	return d.Decoder.Decode(key)
}

// DecoderMaps is a DecoderMap slice.
type DecoderMaps []*DecoderMap

// Decode invokes Decode for each DecoderMap.
func (d DecoderMaps) Decode(key string) string {
	var value string
	for _, dec := range d {
		value = dec.Decode(key)
	}
	return value
}

// AttributeDecoderMaps represents the available DecoderMaps.
var AttributeDecoderMaps = DecoderMaps{
	&DecoderMap{Token: rome.ProcessorToken, Decoder: rome.DecoderMap{Map: rome.Map}},
}
