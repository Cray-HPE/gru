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

package cli

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
)

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

// isDelimiter reports whether the character is a Unicode white space character, or
// a delimiter ",", ";", or "|". The definition of space is set by unicode.IsSpace.
// We avoid dependency on the unicode package, but check validity of the implementation
// in the tests.
func isDelimiter(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		case ',', ';', '|':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

// scanWords is a split function for a Scanner that returns each
// space-separated word of text, with surrounding spaces and other
// delimiters deleted. It will never return an empty string.
func scanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isDelimiter(r) {
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isDelimiter(r) {
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

// ParseHosts checks if os.Stdin is a buffer or not, if it's a
// buffer (e.g. a pipe) then it parses os.Stdin into a []string
// array. If os.Stdin is not a pipe, then the given args are returned
// as is.
func ParseHosts(args []string) []string {
	var newArgs []string
	if isInputFromPipe() {
		scanner := bufio.NewScanner(os.Stdin)
		// Set the split function for the scanning operation.
		scanner.Split(scanWords)
		// Count the words.
		count := 0
		for scanner.Scan() {
			newArgs = append(newArgs, scanner.Text())
			count++
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading input:", err)
		}
		return newArgs
	}
	if len(args) < 1 {
		err := fmt.Errorf("no hosts given")
		if _, exc := fmt.Fprintln(os.Stderr, err); exc != nil {

			panic(exc)
		}
		os.Exit(1)
	}
	return args
}
