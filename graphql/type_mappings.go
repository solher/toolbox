package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/solher/toolbox/sql/types"
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

// MarshalNullString marshals a types.NullString for GqlGen.
func MarshalNullString(t types.NullString) graphql.Marshaler {
	return MarshalString(string(t))
}

// UnmarshalNullString unmarshals a types.NullString for GqlGen.
func UnmarshalNullString(v interface{}) (types.NullString, error) {
	s, err := UnmarshalString(v)
	if err != nil {
		return "", err
	}
	return types.NullString(s), nil
}

// MarshalNullInt64 marshals a types.NullInt64 for GqlGen.
func MarshalNullInt64(t types.NullInt64) graphql.Marshaler {
	return graphql.MarshalInt64(int64(t))
}

// UnmarshalNullInt64 unmarshals a types.NullInt64 for GqlGen.
func UnmarshalNullInt64(v interface{}) (types.NullInt64, error) {
	s, err := graphql.UnmarshalInt64(v)
	if err != nil {
		return 0, err
	}
	return types.NullInt64(s), nil
}

// MarshalNullTimestamp marshals a types.NullTimestamp for GqlGen.
func MarshalNullTimestamp(t types.NullTimestamp) graphql.Marshaler {
	return graphql.MarshalTime(t.Time)
}

// UnmarshalNullTimestamp unmarshals a types.NullTimestamp for GqlGen.
func UnmarshalNullTimestamp(v interface{}) (types.NullTimestamp, error) {
	t, err := graphql.UnmarshalTime(v)
	if err != nil {
		return types.NullTimestamp{}, err
	}
	return types.NullTimestamp{Time: t}, nil
}

// MarshalNullDate marshals a types.NullDate for GqlGen.
func MarshalNullDate(t types.NullDate) graphql.Marshaler {
	return MarshalString(string(t))
}

// UnmarshalNullDate unmarshals a types.NullDate for GqlGen.
func UnmarshalNullDate(v interface{}) (types.NullDate, error) {
	s, err := UnmarshalString(v)
	if err != nil {
		return types.NullDate(""), err
	}
	return types.NullDate(s), nil
}
