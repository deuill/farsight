// Copyright 2016 Alexander Palaistras. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package farsight

import (
	// Standard library.
	"fmt"
	"reflect"

	// Internal packages.
	"github.com/deuill/farsight/parser"
	"github.com/deuill/farsight/source"

	// Pre-defined sources and parsers.
	_ "github.com/deuill/farsight/parser/html"
	_ "github.com/deuill/farsight/source/http"
)

// Fetch data from source pointed to by URI in `src`, and store to arbitrary
// struct pointed to by `dest`.
func Fetch(src string, dest interface{}, kind string) error {
	// Verify destination value type.
	val := reflect.ValueOf(dest)

	if val.Kind() != reflect.Ptr && val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Invalid destination type '%s', expected 'ptr'", val.Kind().String())
	}

	// Fetch data from source defined in `src`.
	buf, err := source.Fetch(src)
	if err != nil {
		return err
	}

	// Parse raw data and return parsed document.
	doc, err := parser.Parse(kind, buf)
	if err != nil {
		return err
	}

	// Populate destination fields from parsed document.
	if err = populateStruct(doc, val.Elem()); err != nil {
		return err
	}

	return nil
}

// Set struct fields sequentially according to their `farsight` tags.
func populateStruct(doc parser.Document, dest reflect.Value) error {
	// Set each struct field in sequence.
	for i := 0; i < dest.NumField(); i++ {
		f := dest.Field(i)
		ft := dest.Type().Field(i)

		// Skip field if `farsight` tag is unset or explicitly ignored.
		attr := ft.Tag.Get("farsight")
		if attr == "" || attr == "-" {
			continue
		}

		// Filter document by tag and set field.
		subdoc, err := doc.Filter(attr)
		if err != nil {
			return err
		}

		if err = setField(subdoc, f); err != nil {
			return err
		}
	}

	return nil
}

// Set struct field for concrete value contained within `doc`.
func setField(doc parser.Document, field reflect.Value) error {
	// Get string value from document.
	val := doc.String()

	// Determine field type and set value, converting if necessary.
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Struct:
		return populateStruct(doc, field)
	default:
		return fmt.Errorf("Unable to set unknown field type '%s'", field.Kind().String())
	}

	return nil
}
