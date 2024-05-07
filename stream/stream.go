package stream

import (
	"math/rand"
	"time"
)

func Filter[E any](s []E, filterFunc func(E) bool) (ret []E) {
	for _, v := range s {
		if filterFunc(v) {
			ret = append(ret, v)
		}
	}
	return
}

func Map[E1, E2 any](s1 []E1, mapFunc func(E1) E2) (s2 []E2) {
	for _, e1 := range s1 {
		s2 = append(s2, mapFunc(e1))
	}
	return s2
}

func MapToAny[E any](s []E) (ret []any) {
	for _, e := range s {
		ret = append(ret, e)
	}
	return ret
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
