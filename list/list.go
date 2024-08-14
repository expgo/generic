package list

func Contains[E comparable](s []E, e E) bool {
	for _, ee := range s {
		if ee == e {
			return true
		}
	}

	return false
}

func ContainsFunc[E any](s []E, e E, matchFunc func(E) bool) bool {
	for _, ee := range s {
		if matchFunc(ee) {
			return true
		}
	}
	return false
}

func Delete[E comparable](s []E, e E) ([]E, bool) {
	for i, ee := range s {
		if ee == e {
			ret := make([]E, 0, len(s))
			ret = append(ret, s...)

			return append(ret[:i], ret[i+1:]...), true
		}
	}

	return s, false
}

func DeleteFunc[E any](s []E, e E, matchFunc func(E) bool) ([]E, bool) {
	for i, ee := range s {
		if matchFunc(ee) {
			ret := make([]E, 0, len(s))
			ret = append(ret, s...)

			return append(ret[:i], ret[i+1:]...), true
		}
	}

	return s, false
}

func Filter[E comparable](s []E, matchFunc func(E) bool) []E {
	ret := make([]E, 0, len(s))

	for _, ee := range s {
		if matchFunc(ee) {
			ret = append(ret, ee)
		}
	}

	return ret
}
