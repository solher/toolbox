package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/solher/toolbox/sql/types"
)

func MarshalNullString(t types.NullString) graphql.Marshaler {
	return graphql.MarshalString(string(t))
}

func UnmarshalNullString(v interface{}) (types.NullString, error) {
	s, err := graphql.UnmarshalString(v)
	if err != nil {
		return "", err
	}
	return types.NullString(s), nil
}
