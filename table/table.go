package table

import (
	"fmt"
	"errors"
	ini "gopkg.in/ini.v1"
	"results/testentry"
)

func LoadTests(filePath string) ([]testentry.TestEntry, error) {
	cfg, err := ini.Load(filePath)
	empty := []testentry.TestEntry(nil)

	// or init TestEntry here, empty, then return if error
	if err != nil {
		return empty, errors.New("Could not load ini file.")
	}

	// Dimensions
	requests := make([]int, 0)
	concurrencyLevels := make([]int, 0)
	routeNames := make([]string, 0)
	paths := make([]string, 0)
	routeDisplayNames := make([]string, 0)

	// TODO deal with cfg == nil
	gen := cfg.Section("general")
	if gen.HasKey("root-url") {
		fmt.Println("Found root url")
		s, _ := gen.GetKey("root-url")
		fmt.Println(s);
	}

	url := gen.Key("root-url").String()

	concurrencyLevels, err = gen.Key("concurrency").StrictInts(",")
	if err != nil {
		return empty, errors.New("Invalid concurrency levels. Must be" +
			" positive integers.")
	}

	requests, err = gen.Key("requests").StrictInts(",")
	if err != nil {
		return empty, errors.New("Invalid number requests. Must be " +
			" positive integers")
	}

	routes := cfg.ChildSections("routes")

	for _, route := range routes {
		name := route.Name()
		display := route.Key("display").String()
		path := route.Key("path").String()

		routeNames = append(routeNames, name)
		paths = append(paths, path)
		routeDisplayNames = append(routeDisplayNames, display)
	}

	testCount := len(concurrencyLevels) * len(requests) * len(paths)
	tests := make([]testentry.TestEntry, testCount)

	// maybe make index a function... could store sizes in a closure.
	var index int
	m := len(requests)
	n := len(paths)

	fmt.Printf("A total of %d tests are needed.\n", testCount)
	for i, _ := range concurrencyLevels {
		for j, _ := range requests {
			for k, _ := range paths {
				index = i*m*n + j*n + k
				tests[index] = testentry.TestEntry{
					Concurrency:	concurrencyLevels[i],
					Requests:		requests[j],
					URL:			url + paths[k], // mb. use full URL
				}
			}
		}
	}

	return tests, nil
}

