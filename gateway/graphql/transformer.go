package graphql

func intFromPointer(v *int) int {
	if v == nil {
		return 0
	}

	return *v
}
