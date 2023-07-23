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

package set

import (
	"fmt"
	"sync"
)

// Async runs an async setting (fn) against a list of hosts. Returns a map of each host with their
// respective set results.
func Async(
	fn func(host string, settings map[string]interface{}) interface{},
	hosts []string,
	data map[string]interface{},
) map[string]any {

	var wg sync.WaitGroup

	sliceLength := len(hosts)
	wg.Add(sliceLength)

	fmt.Printf("Asynchronously querying [%5d] hosts ... \n", len(hosts))

	sm := make(map[string]interface{})

	for _, host := range hosts {

		go func(host string) {

			defer wg.Done()
			sm[host] = fn(host, data)

		}(host)
	}
	wg.Wait()
	return sm
}
