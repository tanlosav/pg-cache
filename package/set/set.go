package set

func SliceToSet(values []string) map[string]struct{} {
	size := len(values)
	set := make(map[string]struct{}, size)

	for i := 0; i < size; i++ {
		set[values[i]] = struct{}{}
	}

	return set
}
