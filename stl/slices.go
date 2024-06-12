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

func MakeUIntSlice(arg ...int) (result []int) {
	switch len(arg) {
	case 1:
		for i := 1; i < arg[0]; i++ {
			result = append(result, i)
		}
	case 2:

		if arg[0] < 0 {
			arg[0] = 1
		}
		for i := arg[0]; i < arg[1]; i++ {
			result = append(result, i)
		}
	case 3:
		if arg[0] < 0 {
			arg[0] = 1
		}
		for i := arg[0]; i < arg[1]; i += arg[2] {
			result = append(result, i)
		}
	}
	return result
}
