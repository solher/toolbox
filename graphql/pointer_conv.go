package graphql

import "github.com/solher/toolbox/graphql/types"

// PointerToStringArray converts a pointer to an array of strings.
func PointerToStringArray(ptr *[]string) []string {
	if ptr != nil {
		return *ptr
	}
	return []string{}
}

// PointerToIDArray converts a pointer to an array of ids.
func PointerToIDArray(ptr *[]types.ID) []types.ID {
	if ptr != nil {
		return *ptr
	}
	return []types.ID{}
}
