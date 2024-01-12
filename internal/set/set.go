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
	"github.com/spf13/viper"
	"sync"
)

// Async calls the given function (fn) asynchronously against each element in a given slice of HTTP/S endpoints. This
// functional also takes a data argument, which is passed along to each async function. Each async call must return an
// interface. Finally, this function returns a map for each result for each endpoint's function call.
func Async(
	fn func(host string, action any) interface{},
	hosts []string,
	data any,
) map[string]any {

	var wg sync.WaitGroup

	sliceLength := len(hosts)
	wg.Add(sliceLength)

	v := viper.GetViper()
	if !v.GetBool("json") {
		fmt.Printf("Asynchronously updating [%5d] hosts ... \n", len(hosts))
	}

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

// AsyncMap is exactly like Async, however the async function is expected to return a map of interfaces.
// TODO: Async and AsyncMap could probably be consolidated with some extra brain juice.
func AsyncMap(
	fn func(host string, actions map[string]interface{}) interface{},
	hosts []string,
	data map[string]interface{},
) map[string]any {

	var wg sync.WaitGroup

	sliceLength := len(hosts)
	wg.Add(sliceLength)

	v := viper.GetViper()
	if !v.GetBool("json") {
		fmt.Printf("Asynchronously updating [%5d] hosts ... \n", len(hosts))
	}

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

// AsyncCall is exactly like Async, however the async function call does not require any arguments be passed to it.
// This is useful when the async function is simply reading or writing value-less data to an endpoint.
func AsyncCall(
	fn func(host string) interface{},
	hosts []string,
) map[string]any {

	var wg sync.WaitGroup

	sliceLength := len(hosts)
	wg.Add(sliceLength)

	v := viper.GetViper()
	if !v.GetBool("json") {
		fmt.Printf("Asynchronously updating [%5d] hosts ... \n", len(hosts))
	}

	sm := make(map[string]interface{})

	for _, host := range hosts {

		go func(host string) {

			defer wg.Done()
			sm[host] = fn(host)

		}(host)
	}
	wg.Wait()
	return sm
}
