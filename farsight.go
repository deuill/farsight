// Copyright 2016 Alexander Palaistras. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package farsight

import (
	// Standard library.
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	// Internal packages.
	"github.com/deuill/farsight/parser"
	"github.com/deuill/farsight/source"

	// Pre-defined sources and parsers.
	_ "github.com/deuill/farsight/parser/html"
	_ "github.com/deuill/farsight/source/http"
)

var (
	regexpValidInt   = `[-+]?[0-9]+`
	regexpValidFloat = `[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?`
)

// Fetch data from source pointed to by URI in `src`, and store to arbitrary
// struct pointed to by `dest`. Data is parsed according to `kind`, and has to
// correspond to a registered parser.
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

// Set struct fields from document, filtered by tags marked by "farsight"
// definitions.
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val == "" {
			field.SetInt(0)
		} else {
			// Truncate string to valid integer.
			val = regexp.MustCompile(regexpValidInt).FindString(val)

			// Parse string and return integer value.
			num, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}

			field.SetInt(num)
		}
	case reflect.Float32, reflect.Float64:
		if val == "" {
			field.SetFloat(0)
		} else {
			// Truncate string to valid floating point number.
			val = regexp.MustCompile(regexpValidFloat).FindString(val)

			// Parse string and return floating point value.
			num, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}

			field.SetFloat(num)
		}
	case reflect.Slice:
		// Decompose document into list and prepare destination slice.
		docs := doc.Slice()
		slice := reflect.MakeSlice(field.Type(), len(docs), cap(docs))

		for i, d := range docs {
			if err := setField(d, slice.Index(i)); err != nil {
				return err
			}
		}

		field.Set(slice)
	case reflect.Struct:
		return populateStruct(doc, field)
	default:
		return fmt.Errorf("Unable to set unknown field type '%s'", field.Kind().String())
	}

	return nil
}
