package utils

import (
	"fmt"
	"strings"
)

// 빈 구조체는 메모리를 사용하지 않음
type void struct{}

var marking void

type Set struct {
	m map[string]void
}

func NewSet() *Set {
	return &Set{make(map[string]void)}
}

// Add, Remove, Contains, Len, Entries, String 메서드를 구현
func (s *Set) Add(value string) {
	s.m[value] = marking
}

func (s *Set) Remove(value string) {
	delete(s.m, value)
}

func (s *Set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

func (s *Set) Len() int {
	return len(s.m)
}

func (s *Set) Entries() []string {
	entries := make([]string, 0, len(s.m))
	for k := range s.m {
		entries = append(entries, k)
	}
	return entries
}

func (s *Set) String() string {
	return fmt.Sprintf("[%s]", strings.Join(s.Entries(), ", "))
}
