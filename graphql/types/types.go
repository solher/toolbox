package types

import (
	"fmt"
	"strconv"
	"time"
)

// ID represents GraphQL's "ID" type.
type ID string

// ImplementsGraphQLType implements the packer.Unmarshaler interface.
func (id ID) ImplementsGraphQLType(name string) bool {
	return name == "ID"
}

// UnmarshalGraphQL implements the packer.Unmarshaler interface.
func (id *ID) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case string:
		*id = ID(input)
		return nil
	default:
		return fmt.Errorf("wrong type")
	}
}

// MarshalJSON implements the exec.marshaler interface.
func (id ID) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, string(id)), nil
}

// Time represents a time (eg. 17:40:38).
type Time string

// ImplementsGraphQLType implements the packer.Unmarshaler interface.
func (t Time) ImplementsGraphQLType(name string) bool {
	return name == "Time"
}

// UnmarshalGraphQL implements the packer.Unmarshaler interface.
func (t *Time) UnmarshalGraphQL(input interface{}) error {
	switch input.(type) {
	default:
		return fmt.Errorf("wrong type")
	}
}

// Date represents a date (eg. 2017-07-22).
type Date string

// ImplementsGraphQLType implements the packer.Unmarshaler interface.
func (d Date) ImplementsGraphQLType(name string) bool {
	return name == "Date"
}

// UnmarshalGraphQL implements the packer.Unmarshaler interface.
func (d *Date) UnmarshalGraphQL(input interface{}) error {
	switch input.(type) {
	default:
		return fmt.Errorf("wrong type")
	}
}

// DateTime represents a timestamp.
type DateTime struct {
	time.Time
}

// ImplementsGraphQLType implements the packer.Unmarshaler interface.
func (t DateTime) ImplementsGraphQLType(name string) bool {
	return name == "DateTime"
}

// UnmarshalGraphQL implements the packer.Unmarshaler interface.
func (t *DateTime) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case time.Time:
		t.Time = input
		return nil
	case string:
		var err error
		t.Time, err = time.Parse(time.RFC3339, input)
		return err
	case int:
		t.Time = time.Unix(int64(input), 0)
		return nil
	case float64:
		t.Time = time.Unix(int64(input), 0)
		return nil
	default:
		return fmt.Errorf("wrong type")
	}
}
