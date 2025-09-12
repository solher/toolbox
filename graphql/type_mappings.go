package graphql

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalString marshals a string for GqlGen.
func MarshalString(t string) graphql.Marshaler {
	return graphql.MarshalString(t)
}

// UnmarshalString unmarshals a string for GqlGen.
func UnmarshalString(v interface{}) (string, error) {
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

// MarshalDate serializes the date as a YYYY-MM-DD string.
func MarshalDate(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}
	// Normalize to midnight UTC
	y, m, d := t.UTC().Date()
	normalized := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(normalized.Format("2006-01-02")))
	})
}

// UnmarshalDate accepts either 'YYYY-MM-DD' or an RFC3339/RFC3339Nano timestamp, and keeps only the date part (midnight UTC).
func UnmarshalDate(v any) (time.Time, error) {
	if s, ok := v.(string); ok {
		// Try strict date first
		if t, err := time.ParseInLocation("2006-01-02", s, time.UTC); err == nil {
			y, m, d := t.Date()
			return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
		}
		// Finally try RFC3339
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			y, m, d := t.Date()
			return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
		}
		// Then try RFC3339Nano
		if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
			y, m, d := t.Date()
			return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
		}
	}
	return time.Time{}, errors.New("date must be 'YYYY-MM-DD' or an RFC3339 formatted string")
}
