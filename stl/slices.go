package stl

func FindElementIndex(s []string, elem string) (int, bool) {
	for i, e := range s {
		if e == elem {
			return i, true
		}
	}
	return -1, false
}

func DeleteElement(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
