package resolver

func safeStringDereference(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}
