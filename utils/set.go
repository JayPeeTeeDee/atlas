package utils

type Set struct {
	vals map[string]struct{}
}

func NewSet() *Set {
	s := &Set{}
	s.vals = make(map[string]struct{})
	return s
}

func (s *Set) Add(value string) {
	s.vals[value] = struct{}{}
}

func (s *Set) AddAll(values ...string) {
	for _, val := range values {
		s.vals[val] = struct{}{}
	}
}

func (s *Set) Remove(value string) {
	delete(s.vals, value)
}

func (s *Set) Contains(value string) bool {
	_, exists := s.vals[value]
	return exists
}

func (s *Set) Keys() []string {
	keys := make([]string, len(s.vals))
	i := 0
	for k := range s.vals {
		keys[i] = k
		i++
	}
	return keys
}

func (s *Set) Union(other *Set) *Set {
	newSet := NewSet()
	for k := range s.vals {
		newSet.Add(k)
	}
	for k := range other.vals {
		newSet.Add(k)
	}
	return newSet
}

func (s *Set) Difference(other *Set) *Set {
	newSet := NewSet()
	for k := range s.vals {
		if !other.Contains(k) {
			newSet.Add(k)
		}
	}
	return newSet
}

func (s *Set) Size() int {
	return len(s.vals)
}
