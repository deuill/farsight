// Copyright 2016 Alexander Palaistras. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package source

import (
	// Standard library.
	"fmt"
	"io"
	"strings"
)

// Fetcher is an interface that wraps the Fetch method.
type Fetcher interface {
	Fetch(src string) (io.Reader, error)
}

// A map of all registered sources.
var sources map[string]Fetcher

// Register a source under a unique name.
func Register(name string, rcvr Fetcher) error {
	if _, exists := sources[name]; exists {
		return fmt.Errorf("Source '%s' already registered, refusing to overwrite", name)
	}

	sources[name] = rcvr
	return nil
}

// Fetch resource, calling the appropriate source handler.
func Fetch(src string) (io.Reader, error) {
	fields := strings.Split(src, ":")
	if len(fields) < 2 {
		return nil, fmt.Errorf("Failed to parse source URL '%s'", src)
	}

	if _, exists := sources[fields[0]]; !exists {
		return nil, fmt.Errorf("Source scheme '%s' does not match a registered source", fields[0])
	}

	return sources[fields[0]].Fetch(src)
}

func init() {
	sources = make(map[string]Fetcher)
}
