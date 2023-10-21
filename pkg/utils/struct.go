package utils

type StringSet struct {
	values map[string]any
}

func (s *StringSet) Contains(value string) bool {
	if s.values == nil {
		return false
	}
	_, ok := s.values[value]
	return ok
}

func NewStringSet(values []string) *StringSet {
	ss := &StringSet{
		values: map[string]any{},
	}
	if values == nil {
		return ss
	}
	for _, value := range values {
		ss.values[value] = true
	}
	return ss
}
