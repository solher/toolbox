package graphql

import (
	"errors"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalString marshals a string for GqlGen.
func MarshalString(v string) graphql.Marshaler {
	return graphql.MarshalString(v)
}

// UnmarshalString unmarshals a string for GqlGen.
func UnmarshalString(v any) (string, error) {
	s, err := graphql.UnmarshalString(v)
	if err != nil {
		return "", err
	}
	// GqlGen sends "null" for null strings. It's weird.
	if s == "null" {
		return "", nil
	}
	return s, nil
}

// MarshalTime serializes the time as a HH:MM:SS string.
func MarshalTime(v string) graphql.Marshaler {
	if v == "" {
		return graphql.Null
	}
	return graphql.MarshalString(v)
}

// UnmarshalTime accepts a 'HH:MM:SS' formatted string.
func UnmarshalTime(v any) (string, error) {
	if s, ok := v.(string); ok {
		if t, err := time.ParseInLocation(time.TimeOnly, s, time.UTC); err == nil {
			return t.Format(time.TimeOnly), nil
		}
	}
	return "", errors.New("time must be 'HH:MM:SS' formatted string")
}

// MarshalDate serializes the date as a YYYY-MM-DD string.
func MarshalDate(v string) graphql.Marshaler {
	if v == "" {
		return graphql.Null
	}
	return graphql.MarshalString(v)
}

// UnmarshalDate accepts a 'YYYY-MM-DD' formatted string.
func UnmarshalDate(v any) (string, error) {
	if s, ok := v.(string); ok {
		if t, err := time.ParseInLocation(time.DateOnly, s, time.UTC); err == nil {
			return t.Format(time.DateOnly), nil
		}
	}
	return "", errors.New("date must be 'YYYY-MM-DD' formatted string")
}

// MarshalTimezone serializes the timezone as a IANA timezone string.
func MarshalTimezone(v time.Location) graphql.Marshaler {
	return graphql.MarshalString(v.String())
}

// UnmarshalTimezone accepts either 'IANA timezone' formatted string.
func UnmarshalTimezone(v any) (time.Location, error) {
	if s, ok := v.(string); ok {
		if l, err := time.LoadLocation(s); err == nil {
			return *l, nil
		}
	}
	return time.Location{}, errors.New("timezone must be a valid IANA timezone string")
}

// MarshalDateTime serializes the datetime as a RFC3339 formatted string.
func MarshalDateTime(v time.Time) graphql.Marshaler {
	if v.IsZero() {
		return graphql.Null
	}
	return graphql.MarshalString(v.Format(time.RFC3339))
}

// UnmarshalDateTime accepts either 'RFC3339' formatted string.
func UnmarshalDateTime(v any) (time.Time, error) {
	if s, ok := v.(string); ok {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("datetime must be a valid RFC3339 formatted string")
}
