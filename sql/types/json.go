// Copyright (c) 2013, Jason Moiron
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONText is a json.RawMessage, which is a []byte underneath.
// Value() validates the json format in the source, and returns an error if
// the json is not valid.  Scan does no validation.  JSONText additionally
// implements `Unmarshal`, which unmarshals the json within to an interface{}
type JSONText json.RawMessage

var emptyJSON = JSONText("{}")

// MarshalJSON returns the *j as the JSON encoding of j.
func (j JSONText) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return emptyJSON, nil
	}
	return j, nil
}

// UnmarshalJSON sets *j to a copy of data
func (j *JSONText) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSONText: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// Value returns j as a value.  This does a validating unmarshal into another
// RawMessage.  If j is invalid json, it returns an error.
func (j JSONText) Value() (driver.Value, error) {
	var m json.RawMessage
	var err = j.Unmarshal(&m)
	if err != nil {
		return []byte{}, err
	}
	return []byte(j), nil
}

// Scan stores the src in *j.  No validation is done.
func (j *JSONText) Scan(src interface{}) error {
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		if len(t) == 0 {
			source = emptyJSON
		} else {
			source = t
		}
	case nil:
		*j = emptyJSON
	default:
		return errors.New("Incompatible type for JSONText")
	}
	*j = JSONText(append((*j)[0:0], source...))
	return nil
}

// Unmarshal unmarshal's the json in j to v, as in json.Unmarshal.
func (j *JSONText) Unmarshal(v interface{}) error {
	if len(*j) == 0 {
		*j = emptyJSON
	}
	return json.Unmarshal([]byte(*j), v)
}

// String supports pretty printing for JSONText types.
func (j JSONText) String() string {
	return string(j)
}

// NullJSONText represents a JSONText that may be null.
// NullJSONText implements the scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullJSONText struct {
	JSONText
	Valid bool // Valid is true if JSONText is not NULL
}

// Scan implements the Scanner interface.
func (n *NullJSONText) Scan(value interface{}) error {
	if value == nil {
		n.JSONText, n.Valid = emptyJSON, false
		return nil
	}
	n.Valid = true
	return n.JSONText.Scan(value)
}

// Value implements the driver Valuer interface.
func (n NullJSONText) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.JSONText.Value()
}
