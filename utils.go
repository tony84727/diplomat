package diplomat

func stringSliceToInterfaceSlice(slice []string) []interface{} {
	pointers := make([]interface{}, len(slice))
	for i, v := range slice {
		s := v
		pointers[i] = &s
	}
	return pointers
}
