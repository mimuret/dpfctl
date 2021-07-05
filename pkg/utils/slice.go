package utils

func StringSliceToInterfaceSlice(src []string) []interface{} {
	r := make([]interface{}, len(src))
	for i, s := range src {
		r[i] = s
	}
	return r
}
