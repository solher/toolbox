package graphql

// PointerToStringArray converts a pointer to an array of string.
func PointerToStringArray(ptr *[]string) []string {
	if ptr != nil {
		return *ptr
	}
	return []string{}
}
