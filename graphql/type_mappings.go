package graphql

import (
	"errors"
	"regexp"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/solher/toolbox/types"
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
func MarshalTime(v types.Time) graphql.Marshaler {
	return graphql.MarshalString(v.String())
}

// UnmarshalTime accepts a 'HH:MM:SS' formatted string.
func UnmarshalTime(v any) (types.Time, error) {
	if s, ok := v.(string); ok {
		if parsed, err := time.ParseInLocation(time.TimeOnly, s, time.UTC); err == nil {
			return types.Time{Time: parsed}, nil
		}
	}
	return types.Time{}, errors.New("time must be 'HH:MM:SS' formatted string")
}

// MarshalDate serializes the date as a YYYY-MM-DD string.
func MarshalDate(v types.Date) graphql.Marshaler {
	return graphql.MarshalString(v.String())
}

// UnmarshalDate accepts a 'YYYY-MM-DD' formatted string.
func UnmarshalDate(v any) (types.Date, error) {
	if s, ok := v.(string); ok {
		if parsed, err := time.ParseInLocation(time.DateOnly, s, time.UTC); err == nil {
			return types.Date{Time: parsed}, nil
		}
	}
	return types.Date{}, errors.New("date must be 'YYYY-MM-DD' formatted string")
}

// MarshalTimeZone serializes the time zone as an IANA time zone string.
func MarshalTimeZone(v types.TimeZone) graphql.Marshaler {
	return graphql.MarshalString(v.String())
}

// UnmarshalTimeZone accepts an IANA time zone string.
func UnmarshalTimeZone(v any) (types.TimeZone, error) {
	if s, ok := v.(string); ok {
		if parsed, err := time.LoadLocation(s); err == nil {
			return types.TimeZone{Location: *parsed}, nil
		}
	}
	return types.TimeZone{}, errors.New("timezone must be a valid IANA time zone string")
}

// MarshalDateTime serializes the datetime as a RFC3339 formatted string.
func MarshalDateTime(v time.Time) graphql.Marshaler {
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

var countryRegex = regexp.MustCompile("^[A-Z]{2}$")

// MarshalCountry serializes the country as an ISO 3166-1 alpha-2 code.
func MarshalCountry(v string) graphql.Marshaler {
	return graphql.MarshalString(v)
}

// UnmarshalCountry ensures the country is an ISO 3166-1 alpha-2 code.
func UnmarshalCountry(v any) (string, error) {
	if s, ok := v.(string); ok {
		if countryRegex.MatchString(s) {
			return s, nil
		}
	}
	return "", errors.New("country must be a valid ISO 3166-1 alpha-2 code")
}

var languageRegex = regexp.MustCompile("^[a-z]{2}(-[0-9A-Z]+)?$")

// MarshalLanguage serializes the language as an IETF language tag.
func MarshalLanguage(v string) graphql.Marshaler {
	return graphql.MarshalString(v)
}

// UnmarshalLanguage ensures the language is an IETF language tag.
func UnmarshalLanguage(v any) (string, error) {
	if s, ok := v.(string); ok {
		if languageRegex.MatchString(s) {
			return s, nil
		}
	}
	return "", errors.New("language must be a valid IETF language tag")
}
