package stream

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

var (
	/*
		Used by FlatMap, Map and ToMap to control the iterator
	*/
	// ErrContinue is an error that signals to continue the operation.
	ErrContinue = errors.New("continue")
	// ErrBreak is an error that signals to break the operation.
	ErrBreak = errors.New("break")
)

// Stream is a generic type representing a stream of elems of type T.
// It contains a slice of elems of type T.
type Stream[T any] struct {
	elems []T
	err   error
}

// Of creates a new Stream from the provided elements.
// It takes a slice of elements as input and returns a Stream
// with the same elements.
//
// Example:
//
//	s := stream.Of([]int{1, 2, 3, 4, 5})
//
// Parameters:
//   - elems: The elements to create the Stream from.
//
// Returns:
//   - A Stream containing the provided elements.
func Of[T any](elems []T) Stream[T] {
	return Stream[T]{elems: elems}
}

// Append appends the given values to the stream tail.
// The original stream is not modified.
//
// The elements are appended in the order they are supplied.
// The appended values can be of any type specified by T in the stream declaration.
func (s Stream[T]) Append(values ...T) Stream[T] {
	if s.err != nil {
		return Stream[T]{err: s.err}
	}

	s.elems = append(s.elems, values...)
	return s
}

// Filter applies the filterFunc function to each element in the stream and returns a new stream that contains only
// the elements for which the filter function returns true.
// The original stream is not modified.
// The filter function takes an argument of type T and returns a boolean value indicating whether the element should be included in the filtered stream.
// If the filter function returns an error, the filtering process is stopped and the resulting stream will have its err field set to the error value.
// The filtered elements are appended to the result stream in the order they are encountered.
func (s Stream[T]) Filter(filterFunc func(T) (bool, error)) Stream[T] {
	if s.err != nil {
		return Stream[T]{err: s.err}
	}

	var result Stream[T]

	for i, v := range s.elems {
		ok, err := filterFunc(v)
		if err != nil {
			result.err = fmt.Errorf("stream filter elems[%d] with err: %v", i, err)
			return result
		}

		if ok {
			result.elems = append(result.elems, v)
		}
	}

	return result
}

// Shuffle randomly rearranges the elements in the stream.
// It creates a new stream with the same elements as the original stream, but in a random order.
// The original stream remains unchanged.
// The shuffle algorithm used is the Fisher-Yates shuffle.
// The seed for the random number generator is set using the current time.
// Example usage:
//
//	stream := Of([]int{1, 2, 3, 4, 5})
//	shuffled := stream.Shuffle()
//	shuffledElems := shuffled.MustToSlice()
//	fmt.Println(shuffledElems)  // Output: [4 3 1 2 5]
func (s Stream[T]) Shuffle() Stream[T] {
	if s.err != nil {
		return Stream[T]{err: s.err}
	}

	if len(s.elems) == 0 {
		return s
	}

	//Create a new Stream and copy the data from the original Stream over
	newStream := Stream[T]{elems: append([]T(nil), s.elems...)}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < r.Intn(3)+3; i++ {
		for n := len(newStream.elems); n > 0; n-- {
			randIndex := r.Intn(n)
			newStream.elems[n-1], newStream.elems[randIndex] = newStream.elems[randIndex], newStream.elems[n-1]
		}
	}

	return newStream
}

// Distinct returns a new stream that contains only distinct elements from the original stream.
// The distinctness of elements is determined by the provided equalFunc function.
//
// The original stream is not modified.
// The elements are iterated in the order they appear in the original stream.
// The equalFunc function should take two elements, preItem and nextItem, and return true if they are equal, and false otherwise.
// It can also return an error if there is an error during the comparison.
//
// Example usage:
//
//	equalFunc := func(preItem, nextItem T) (bool, error) {
//	  // implementation of equality comparison logic
//	}
//	distinctStream := stream.Distinct(equalFunc)
func (s Stream[T]) Distinct(equalFunc func(preItem, nextItem T) (bool, error)) Stream[T] {
	if s.err != nil {
		return Stream[T]{err: s.err}
	}

	if len(s.elems) == 0 {
		return s
	}

	var result Stream[T]
	result.elems = append(result.elems, s.elems[0])

	for _, newItem := range s.elems[1:] {
		unique := true
		for _, existingItem := range result.elems {
			equal, err := equalFunc(existingItem, newItem)
			if err != nil {
				result.err = err
				return result
			}
			if equal {
				unique = false
				break
			}
		}
		if unique {
			result.elems = append(result.elems, newItem)
		}
	}

	return result
}

// Reverse reverses the order of elements in the stream.
// If the stream contains an error, it returns a new Stream with the same error.
// The original stream is not modified.
// It uses a two-pointer technique to swap elements starting from both ends of the stream until they meet in the middle.
func (s Stream[T]) Reverse() Stream[T] {
	if s.err != nil {
		return Stream[T]{err: s.err}
	}

	if len(s.elems) == 0 {
		return s
	}

	//Create a new Stream and copy the data from the original Stream over
	newStream := Stream[T]{elems: append([]T(nil), s.elems...)}

	for i, j := 0, len(newStream.elems)-1; i < j; i, j = i+1, j-1 {
		newStream.elems[i], newStream.elems[j] = newStream.elems[j], newStream.elems[i]
	}

	return newStream
}

// Sort sorts the elements in the stream in ascending order according to the compareFunc.
// It uses the sort.Slice function to perform the sorting.
// The compareFunc should take two elements of type T and return an integer value.
// If the value is less than 0, the first element is considered smaller than the second element.
// If the value is equal to 0, the elements are considered equal.
// If the value is greater than 0, the first element is considered greater than the second element.
// The original stream is modified in place.
// The sorted stream is returned.
// Example Usage:
//
//	compare := func(a, b int) int {
//	    return a - b
//	}
//	stream := Of([]int{3, 1, 2}).Sort(compare)
//	result := stream.MustToSlice() // [1, 2, 3]
func (s Stream[T]) Sort(compareFunc func(x, y T) int) Stream[T] {
	if s.err != nil {
		return Stream[T]{err: s.err}
	}

	if len(s.elems) == 0 {
		return s
	}

	//Create a new Stream and copy the data from the original Stream over
	newStream := Stream[T]{elems: append([]T(nil), s.elems...)}

	sort.Slice(newStream.elems, func(i, j int) bool {
		return compareFunc(newStream.elems[i], newStream.elems[j]) < 0
	})

	return newStream
}

// Limit returns a new Stream containing at most `n` elements from the current Stream.
// If `n` is negative, it is set to 0.
// If `n` is greater than the number of elements in the current Stream, it is set to the number of elements.
// The order of the elements in the new Stream is the same as in the current Stream.
// The new Stream is returned as a pointer to Stream[T].
//
// Example usage:
//
//	s := NewStream([]int{1, 2, 3, 4, 5})
//	limited := s.Limit(3)
//	limited.ToSlice() // returns [1, 2, 3]
func (s Stream[T]) Limit(n int) Stream[T] {
	if s.err != nil {
		return Stream[T]{err: s.err}
	}

	if n < 0 {
		n = 0
	} else if n > len(s.elems) {
		n = len(s.elems)
	}

	return Of(s.elems[:n])
}

// Skip skips the first `n` elements in the stream and returns a new stream without those elements.
// If `n` is negative, Skip behaves as if `n` is 0.
// If `n` is greater than the number of elements in the stream, Skip behaves as if `n` is equal to the number of elements in the stream.
// The original stream is not modified.
// The elements of the new stream are in the same order as in the original stream, starting from the `n+1`th element.
// The new stream is returned as a pointer to Stream[T].
// Example usage:
//
//	s := Of([]int{1, 2, 3, 4, 5})
//	newStream := s.Skip(2)
//	fmt.Println(newStream.ToSlice()) // Output: [3 4 5]
//	fmt.Println(s.ToSlice()) // Output: [1 2 3 4 5]
func (s Stream[T]) Skip(n int) Stream[T] {
	if s.err != nil {
		return Stream[T]{err: s.err}
	}

	if n < 0 {
		n = 0
	} else if n > len(s.elems) {
		n = len(s.elems)
	}

	return Of(s.elems[n:])
}

// Map applies the given map function to each element in the stream and returns a new stream containing the results.
// The original stream is not modified.
// The map function takes an element of type T as input and returns an element of type T.
// The returned stream will have the same type T as the original stream.
// If the original stream has an error, a new stream with the same error will be returned.
// If the map function encounters an error during the mapping operation, the error will be propagated and the mapping process will be stopped.
func (s Stream[T]) Map(mapFunc func(T) (T, error)) Stream[T] {
	if s.err != nil {
		return Stream[T]{err: s.err}
	}

	result := Stream[T]{}

	for i, elem := range s.elems {
		mappedValue, err := mapFunc(elem)
		if err != nil {
			result.err = fmt.Errorf("stream map elems[%d] with err: %v", i, err)
			return result
		}

		result.elems = append(result.elems, mappedValue)
	}

	return result
}

// ToSlice returns a slice containing all the elements of the stream.
// The original stream is not modified.
// The elements in the returned slice are in the same order as in the original stream.
// The returned slice has the type []T, where T is the type of elements in the stream.
// Example usage:
//
//	stream := Of([]int{1, 2, 3})
//	result, _ := stream.ToSlice() // result is []int{1, 2, 3}
func (s Stream[T]) ToSlice() ([]T, error) {
	return s.elems, s.err
}

// ToAny converts the elements of the stream to the `any` type and returns them as a slice.
// It creates a new slice and appends the converted elements of the stream to it.
// The original stream is not modified.
// The elements in the resulting slice follow the same order as in the original stream.
// The resulting slice is returned as a value of type `[]any`.
func (s Stream[T]) ToAny() ([]any, error) {
	if s.err != nil {
		return nil, s.err
	}

	var result []any

	for _, v := range s.elems {
		result = append(result, any(v))
	}

	return result, nil
}

// Contains checks if the stream contains a specific value.
// It returns true if the value is found in the stream,
// otherwise it returns false. The value to search for
// is specified by the parameter 'value' which can be of
// any type specified by T in the stream declaration.
// The equalFunc parameter is a function that takes two
// arguments of type T and returns a boolean value and an
// error. This function is used to compare each element in
// the stream with the search value. If an error occurs during
// the comparison, it will be returned.
// The function iterates over each element in the stream and
// performs the comparison using equalFunc. If a match is found,
// it returns true. If no match is found, it returns false.
// If the stream contains any error, it will be returned along
// with false to indicate that the search operation was not
// successful. Otherwise, if the search was successful, it will
// return true and nil.
// Example usage:
//
//	stream := Stream["string"]{elems: []string{"a", "b", "c"}}
//	contains, err := stream.Contains("b", func(x, y string) (bool, error) {
//	  return x == y, nil
//	})
//	if err != nil {
//	  fmt.Println("An error occurred:", err)
//	}
//	fmt.Println("Contains:", contains)
//	// Output: Contains: true
func (s Stream[T]) Contains(value T, equalFunc func(x, y T) (bool, error)) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for _, v := range s.elems {
		ok, err := equalFunc(v, value)
		if err != nil {
			return false, err
		}

		if ok {
			return true, nil
		}
	}

	return false, nil
}

// AllMatch returns true if all elements in the stream satisfy the given matchFunc function.
//
// It iterates through each element in the stream and applies the matchFunc function to determine if the element satisfies the condition.
// If any element fails the condition, the function immediately returns false.
// If all elements pass the condition, the function returns true.
//
// The original stream is not modified.
// The matchFunc function should return true for elements that satisfy the condition, and false for elements that do not.
// The stream is a pointer to Stream[T] type.
//
// Example usage:
//
//	stream := Of([]T{1, 2, 3, 4, 5})
//	result := stream.AllMatch(func(elem T) bool {
//	  return elem > 0
//	})
//	// result is true, since all elements in the stream are greater than 0
//
// Note: The elements of the stream should be of the same type as the type specified for Stream[T].
// For example, if the Stream[T] is created with Stream[int], the elements should be of type int.
// The behavior of the method is undefined if this condition is violated.
func (s Stream[T]) AllMatch(matchFunc func(T) (bool, error)) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for _, elem := range s.elems {
		match, err := matchFunc(elem)
		if err != nil {
			return false, err
		}

		if !match {
			return false, nil
		}
	}

	return true, nil
}

// AnyMatch checks if any element in the stream satisfies the given matchFunc.
// It iterates over each element in the stream and applies the matchFunc function to it.
// If the matchFunc returns true for any element, the method returns true.
// If the matchFunc returns false for all elements, the method returns false.
// The original stream is not modified.
// The matchFunc function should return true for elements that satisfy the condition and false otherwise.
// Returns true if any element in the stream satisfies the matchFunc, false otherwise.
func (s Stream[T]) AnyMatch(matchFunc func(T) (bool, error)) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for _, elem := range s.elems {
		match, err := matchFunc(elem)
		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}

	return false, nil
}

// FindFirst returns the first element in the stream that satisfies the equalFunc function.
// It iterates through each element in the stream and applies the equalFunc function to determine if it should be kept.
// The equalFunc function should return true for elements that should be returned as the first element, and false for elements that should be skipped.
// If an error occurs while applying the equalFunc function, that error is returned along with the default value for type T.
// If no matching element is found, an error is returned with the default value for type T.
// The original stream is not modified.
// FindFirst returns the first matching element as type T and an error.
// Example:
//
//	stream := Stream[int]{1, 2, 3, 4, 5}
//	equalFunc := func(n int) (bool, error) {
//	    return n > 3, nil
//	}
//	first, err := stream.FindFirst(equalFunc)
//	// first = 4, err = nil
//
// Note: Replace "T" with the actual type used in the implementation.
func (s Stream[T]) FindFirst(equalFunc func(T) (bool, error)) (t T, err error) {
	if s.err != nil {
		return t, s.err
	}

	ok := false
	for _, elem := range s.elems {
		ok, err = equalFunc(elem)
		if err != nil {
			return t, err
		}
		if ok {
			return elem, nil
		}
	}

	return t, errors.New("no matching element found")
}

// Max finds the maximum value in the stream using a compare function.
// It takes a compare function as input, which is responsible for comparing two elements and returning an integer value.
// The compare function should return a positive integer if the second element is greater than the first,
// a negative integer if the second element is smaller than the first,
// and zero if the two elements are equal.
// If the stream has an error, the function returns the zero value of type T and the error from the stream.
// If the stream is empty, it returns the zero value of type T and an error indicating that the stream is empty.
// If the stream is not empty, the function iterates through the elements using a loop,
// comparing each element with the current maximum value and updating it if a larger value is found.
// The final maximum value is returned along with a nil error.
func (s Stream[T]) Max(compareFunc func(x, y T) (int, error)) (t T, err error) {
	if s.err != nil {
		return t, s.err
	}

	if len(s.elems) == 0 {
		return t, errors.New("elems is empty")
	}

	maxValue := s.elems[0]
	compared := 0
	for _, elem := range s.elems[1:] {
		compared, err = compareFunc(elem, maxValue)
		if err != nil {
			return t, err
		}

		if compared > 0 {
			maxValue = elem
		}
	}

	return maxValue, nil
}

// Min returns the minimum value in the stream, based on the provided comparison function.
// The original stream is not modified.
//
// The comparison function (compareFunc) is used to determine the order of the elements.
// It should accept two elements of type T and return an integer value, where:
//   - a negative value indicates that the first element is smaller,
//   - zero indicates equality, and
//   - a positive value indicates that the first element is greater.
//
// If there is an error during the comparison, it should be returned as the second return value.
//
// If the stream is empty, an error is returned indicating that the elements are empty.
func (s Stream[T]) Min(compareFunc func(x, y T) (int, error)) (t T, err error) {
	if s.err != nil {
		return t, s.err
	}

	if len(s.elems) == 0 {
		return t, errors.New("elems is empty")
	}

	minValue := s.elems[0]
	compared := 0
	for _, elem := range s.elems[1:] {
		compared, err = compareFunc(elem, minValue)
		if err != nil {
			return t, err
		}

		if compared < 0 {
			minValue = elem
		}
	}

	return minValue, nil
}

// First returns the first element of the stream.
// If the stream has an error, it returns the default value for type T and the error.
// If the stream is empty, it returns the default value for type T and an error indicating that the stream is empty.
func (s Stream[T]) First() (t T, err error) {
	if s.err != nil {
		return t, s.err
	}

	if len(s.elems) == 0 {
		return t, errors.New("stream is empty")
	}

	return s.elems[0], nil
}

func (s Stream[T]) _reduce(initItem T, beginItem int, accumulator func(preItem, nextItem T) (T, error)) (t T, err error) {
	result := initItem
	for i := beginItem; i < len(s.elems); i++ {
		result, err = accumulator(result, s.elems[i])
		if err != nil {
			return t, err
		}
	}

	return result, nil
}

// Reduce applies the accumulator function to each element in the stream,
// starting with the initial value and the second element.
// The result is passed as the first argument of the accumulator function,
// along with the next element of the stream as the second argument.
// This process continues until all elements in the stream have been processed.
// The accumulator function should return the accumulated value and an error.
// If the accumulator function returns an error at any point, Reduce will stop processing and return that error.
// If the stream is empty, Reduce will return the initial value and a nil error.
// Reduce returns the accumulated value and an error as a tuple.
// Without duplicating the example above, document the following code:
func (s Stream[T]) Reduce(accumulator func(preItem, nextItem T) (T, error)) (t T, err error) {
	if s.err != nil {
		return t, s.err
	}

	if len(s.elems) == 0 {
		return t, nil
	}

	return s._reduce(s.elems[0], 1, accumulator)
}

// ReduceWithInit reduces the stream by applying the accumulator function to each element,
// starting with the initial value initItem. The accumulator function takes the previous
// accumulated value and the next element of the stream, and returns the new accumulated value.
// The result of the reduction is the final accumulated value.
// If the stream is empty, the initial value initItem is returned.
// The accumulator function can also return an error, which will be propagated and
// cause the reduction to stop and return the error immediately.
// The initial value initItem can be of any type T, and the return value of the accumulator
// function must be of type T as well.
// The original stream remains unchanged.
// The reduction is performed in the order of the elements in the stream.
// The result of the reduction and any error encountered are returned as a tuple.
// If an error is encountered during the reduction, the value of the result is undefined.
func (s Stream[T]) ReduceWithInit(initItem T, accumulator func(preItem, nextItem T) (T, error)) (t T, err error) {
	if s.err != nil {
		return t, s.err
	}

	if len(s.elems) == 0 {
		return t, nil
	}

	return s._reduce(initItem, 0, accumulator)
}

// Range iterates over each element in the stream and applies the forEachFun function to it.
// If the forEachFun function returns false for any element, the iteration is stopped.
// The forEachFun function should return true for elements that need to be processed, and false for elements that can be skipped.
// This method does not modify the original stream.
// The elements are iterated in the same order as in the stream.
// This method does not return any value.
func (s Stream[T]) Range(forEachFun func(T) error) error {
	if s.err != nil {
		return s.err
	}

	for _, elem := range s.elems {
		if err := forEachFun(elem); err != nil {
			return err
		}
	}

	return nil
}

// Size returns the number of elements in the stream.
// It calculates and returns the length of s.elems.
// The count includes all elements in the stream, regardless of any filters applied.
// The returned value is an integer representing the size of the stream.
func (s Stream[T]) Size() int {
	return len(s.elems)
}

// Err returns the error associated with the stream.
// It retrieves the error value stored in the 'err' field of the Stream struct.
// This method can be used to check if an error occurred during stream processing.
// If no error occurred, it returns nil.
// The error value is returned as an instance of the 'error' interface.
// Example usage:
//
//	stream := &Stream[T]{}
//	err := stream.Err()
//	if err != nil {
//	    fmt.Println("An error occurred:", err.Error())
//	}
func (s Stream[T]) Err() error {
	return s.err
}

// GroupBy groups the elements of the input stream based on the provided key function.
// It returns a map where each key corresponds to a group, and the value is a stream
// containing the elements that belong to that group.
// The key function determines the grouping criteria by extracting a key from each element.
// If two elements have the same key, they will belong to the same group.
//
// Example:
//
//	s := Of([]int{1, 2, 3, 4, 5})
//	groups := GroupBy(s, func(num int) string {
//	    if num%2 == 0 {
//	        return "even"
//	    }
//	    return "odd"
//	})
//
//	// The resulting groups map will be:
//	// {
//	//    "even": {elems: [2, 4]},
//	//    "odd": {elems: [1, 3, 5]},
//	// }
//
// If the input stream is empty, the result will be an empty map.
//
// Parameters:
// - s: The input stream to group.
// - getKey: A function that extracts the key from each element in the stream.
//
// Returns:
//   - A map where each key corresponds to a group, and the value is a stream
//     containing the elements that belong to that group.
func GroupBy[T any, K comparable](s Stream[T], getKey func(T) K) map[K]Stream[T] {
	result := make(map[K]Stream[T])
	if s.err != nil {
		return result
	}

	for _, v := range s.elems {
		key := getKey(v)
		if _, ok := result[key]; !ok {
			result[key] = Of([]T{v})
		} else {
			result[key] = result[key].Append(v)
		}
	}

	return result
}

// ToMap converts the elements of the input stream into a map using the provided map function.
// The map function takes an element of the input stream and returns a key-value pair and an optional error.
// If the input stream has an error, it is returned as it is.
// The resulting map is returned along with a potential error.
// If the map function returns an error, the conversion stops and the error is returned immediately.
// The keys and values in the map are of types Key and value, respectively.
func ToMap[In any, Key comparable, Value any](in Stream[In], mapFunc func(In) (Key, Value, error)) (map[Key]Value, error) {
	if in.err != nil {
		return nil, in.err
	}

	result := make(map[Key]Value)

	for _, v := range in.elems {
		key, value, err := mapFunc(v)
		if err != nil {
			if errors.Is(err, ErrContinue) {
				continue
			} else if errors.Is(err, ErrBreak) {
				return result, nil
			} else {
				return nil, err
			}
		}
		result[key] = value
	}
	return result, nil
}

// Map applies the provided function `f` to each element of the input stream `s`
// and returns a new stream containing the resulting elements. If an error occurs during
// the mapping process, the resulting stream will have the corresponding error value.
//
// Example:
//
//	str := Of([]int{1, 2, 3, 4})
//	double := func(i int) (int, error) {
//	    return i * 2, nil
//	}
//	doubledStr := Map(str, double)
//	fmt.Println(doubledStr.ToSlice()) // Output: [2, 4, 6, 8]
//
//	str2 := Of([]string{"hello", "world"})
//	length := func(s string) (int, error) {
//	    return len(s), nil
//	}
//	lengthStr := Map(str2, length)
//	fmt.Println(lengthStr.ToSlice()) // Output: [5, 5]
func Map[In any, Out any](s Stream[In], f func(In) (Out, error)) Stream[Out] {
	var result Stream[Out]
	if s.err != nil {
		result.err = s.err
		return result
	}

	for _, v := range s.elems {
		elem, err := f(v)
		if err != nil {
			if errors.Is(err, ErrContinue) {
				continue
			} else if errors.Is(err, ErrBreak) {
				return result
			} else {
				result.err = err
				return result
			}
		}

		result.elems = append(result.elems, elem)
	}

	return result
}

// FlatMap takes a Stream `in` and a `flatMap` function and applies the `flatMap` function to each element in the Stream `in`.
// It returns a new Stream with the concatenated elements from all the resulting Streams.
// If the Stream `in` has an error, the error is propagated to the result Stream.
// If any resulting Stream from the `flatMap` function has an error, the error is also propagated to the result Stream.
func FlatMap[In any, Out any](in Stream[In], flatMap func(In) Stream[Out]) Stream[Out] {
	var result Stream[Out]
	if in.err != nil {
		result.err = in.err
		return result
	}

	for _, v := range in.elems {
		stream := flatMap(v)
		if stream.err != nil {
			if errors.Is(stream.err, ErrContinue) {
				continue
			} else if errors.Is(stream.err, ErrBreak) {
				return result
			} else {
				result.err = stream.err
				return result
			}
		}

		result.elems = append(result.elems, stream.elems...)
	}

	return result
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
