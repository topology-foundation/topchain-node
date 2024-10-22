package types

// Set implementation from
// https://medium.com/@dpinoagustin/implementing-a-set-type-in-golang-437174891d99
// Set is the set type
type Set[T comparable] map[T]struct{}

// Add returns the current Set by adding a new element.
func (s Set[T]) Add(e T) Set[T] {
	s[e] = struct{}{}
	return s
}

// Remove returns the current Set by removing an element.
func (s Set[T]) Remove(e T) Set[T] {
	delete(s, e)
	return s
}

// Has returns true when the element is in the Set. Otherwise, returns false.
func (s Set[T]) Has(e T) bool {
	_, ok := s[e]
	return ok
}

// From returns a new Set with the elements provided.
func SetFrom[T comparable](elms ...T) Set[T] {
	s := make(Set[T], len(elms))
	for _, e := range elms {
		s.Add(e)
	}
	return s
}
