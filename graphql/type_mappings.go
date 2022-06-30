package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/solher/toolbox/sql/types"
)

// MarshalNullString marshals a types.NullString for GqlGen.
func MarshalNullString(t types.NullString) graphql.Marshaler {
	return graphql.MarshalString(string(t))
}

// UnmarshalNullString unmarshals a types.NullString for GqlGen.
func UnmarshalNullString(v interface{}) (types.NullString, error) {
	s, err := graphql.UnmarshalString(v)
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
