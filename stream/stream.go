package stream

import (
	"math/rand"
	"time"
)

func Limit[E any](s []E, n int) []E {
	if n < 0 {
		n = 0
	} else if n > len(s) {
		n = len(s)
	}
	return s[:n]
}

func Skip[E any](s []E, n int) []E {
	if n < 0 {
		n = 0
	} else if n > len(s) {
		n = len(s)
	}
	return s[n:]
}

func Filter[E any](s []E, filterFunc func(E) bool) (ret []E) {
	for _, v := range s {
		if filterFunc(v) {
			ret = append(ret, v)
		}
	}
	return
}

func Shuffle[E any](s []E) (ret []E) {
	if len(s) == 0 {
		return
	}

	//Create a new Stream and copy the data from the original Stream over
	ret = append([]E(nil), s...)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < r.Intn(3)+3; i++ {
		for n := len(ret); n > 0; n-- {
			randIndex := r.Intn(n)
			ret[n-1], ret[randIndex] = ret[randIndex], ret[n-1]
		}
	}

	return ret
}

func Distinct[E comparable](s []E) []E {
	seen := make(map[E]bool)
	ret := make([]E, 0, len(s))
	for _, v := range s {
		if !seen[v] {
			ret = append(ret, v)
			seen[v] = true
		}
	}
	return ret
}

func DistinctFunc[E any](s []E, matchFunc func(preItem, nextItem E) bool) []E {
	ret := make([]E, 0, len(s))

	if len(s) == 0 {
		return ret
	}

	ret = append(ret, s[0])

	for _, newItem := range s[1:] {
		unique := true
		for _, existingItem := range ret {
			if matchFunc(existingItem, newItem) {
				unique = false
				break
			}
		}
		if unique {
			ret = append(ret, newItem)
		}
	}

	return ret
}

func AllMatch[E comparable](s []E, e E) bool {
	for _, elem := range s {
		if elem != e {
			return false
		}
	}
	return true
}

func AllMatchFunc[E any](s []E, matchFunc func(E) bool) bool {
	for _, elem := range s {
		if !matchFunc(elem) {
			return false
		}
	}
	return true
}

func AnyMatch[E comparable](s []E, e E) bool {
	for _, elem := range s {
		if elem == e {
			return true
		}
	}
	return false
}

func AnyMatchFunc[E any](s []E, matchFunc func(E) bool) bool {
	for _, elem := range s {
		if matchFunc(elem) {
			return true
		}
	}
	return false
}

func ToAny[E any](s []E) (ret []any) {
	for _, e := range s {
		ret = append(ret, e)
	}
	return ret
}

func MustMap[E1, E2 any](s1 []E1, mapFunc func(E1) E2) (s2 []E2) {
	for _, e1 := range s1 {
		s2 = append(s2, mapFunc(e1))
	}
	return s2
}

func Map[E1, E2 any](s1 []E1, mapFunc func(E1) (E2, error)) (s2 []E2, e error) {
	for _, e1 := range s1 {
		e2, err := mapFunc(e1)
		if err != nil {
			return nil, err
		}
		s2 = append(s2, e2)
	}
	return s2, nil
}

func GroupBy[E any, K comparable](s []E, getKey func(E) K) map[K][]E {
	result := make(map[K][]E)

	for _, v := range s {
		key := getKey(v)
		if _, ok := result[key]; !ok {
			result[key] = []E{v}
		} else {
			result[key] = append(result[key], v)
		}
	}

	return result
}
