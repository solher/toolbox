package graphql

import (
	"fmt"
	"strconv"

	"github.com/solher/toolbox/graphql/types"
)

// ToGlobalID takes a type name and an ID specific to that type name, and returns a
// "global ID" that is unique among all types. ObjCode is two characters.
func ToGlobalID(objCode string, id uint64) types.ID {
	str := fmt.Sprintf("%s%d", objCode, id)
	// encStr := base64.StdEncoding.EncodeToString([]byte(str))
	return types.ID(str)
}

// FromGlobalID takes the "global ID" created by toGlobalID, and returns the type name and ID
// used to create it. ObjCode is two characters.
func FromGlobalID(globalID types.ID) (objCode string, id uint64) {
	// idStr, err := base64.StdEncoding.DecodeString(string(globalID))
	if len(globalID) < 3 {
		return "", 0
	}
	id, err := strconv.ParseUint(string(globalID[2:]), 10, 64)
	if err != nil {
		return "", 0
	}
	return string(globalID[:2]), id
}
