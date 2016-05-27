# Farsight - Fetch, filter, and store arbitrary data

[![API Documentation][godoc-svg]][godoc-url] [![MIT License][license-svg]][license-url]

Farsight facilitates the fetching and transformation of data from arbitrary sources into pre-defined structures, which can be further processed and serialised into other formats (such as JSON, YAML etc.).

A large amount of inspiration comes from the [GoStruct](https://github.com/bfontaine/gostruct) project.

## Usage

Farsight is very simple to use, and only exposes a single public method, `Fetch`. For example, this is the full package for a single-page HTML scraper:

```go
package main

import (
	"fmt"
	"github.com/deuill/farsight"
)

var url = "https://deuill.org"

type Post struct {
	Title string `farsight:"h2"`
	Text  string `farsight:".post-text"`
}

type Data struct {
	Intro string `farsight:".post-text:first-of-type"`
	Posts []Post `farsight:".post-summary"`
}

func main() {
	data := &Data{}
	if err := farsight.Fetch(url, data, "html"); err != nil {
		panic("Failed to fetch URL")
	}

	fmt.Println(data.Posts[0].Title) // Returns the first post's title.
}
```

Calling the `Fetch` method on the `Data` type fills each eligible field with the correct data, as matched by the specified CSS selectors.

## Overview

Data sourcing is handled via generic `source` types, that correspond to URIs passed to `farsight.Fetch` and allow for transparent use of different types of sources (such as local files, HTTP endpoints etc.).

Data transformation is handled via `parser` types, that rely on specific struct tag fields in order to fill the destination structures with the correct data. Thus, there is a direct relationship between the type of data the user expects to query against, and the `parser` type used.

## License

Farsight is licensed under the MIT license, the terms of which can be found in the included LICENSE file.


[godoc-url]: https://godoc.org/github.com/deuill/farsight
[godoc-svg]: https://godoc.org/github.com/deuill/farsight?status.svg

[license-url]: https://github.com/deuill/farsight/blob/master/LICENSE
[license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
