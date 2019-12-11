package main

import (
	"fmt"
	"bytes"
)

func main() {
	var s, c IntSet
	s.Add(1)
	s.Add(198)
	s.Add(6954)

	fmt.Printf("s=%s\n", &s)

	c.Add(10)
	c.Add(198)
	c.Add(9999)
	c.Add(44)

	fmt.Printf("c=%s\n", &c)

	s.UnionWith(&c)
	fmt.Printf("c|s=%s\n", &s)

	fmt.Printf("len(s)=%d\n", s.Len())

	s.Remove(10)
	fmt.Printf("remove 10, s=%s\n", &s)

	t := s.Copy()
	s.Clear()
	fmt.Printf("clear s=%s\n", &s)
	fmt.Printf("copy s=%s\n", t)

	s.AddAll(1, 2, 5, 8)
	fmt.Printf("s.AddAll=%s\n", &s)

}

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) bool {
	word, bit := x / 64, uint(x % 64)
	return word < len(s.words) && s.words[word] & (1 << bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x / 64, uint(x % 64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) UnionWith( t *IntSet) {
	for i, word := range t.words {
		if i < len(s.words) {
			s.words[i] |= word	
		}else {
			s.words = append(s.words, word)
		}
	}
}

func (s *IntSet) Len() int {
	count := 0
	for _, word := range s.words {
		for word > 0 {
			word &= word - 1
			count++
		}
	}
	return count
}

func (s *IntSet) Remove(x int) {
	word, bit := x / 64, uint(x % 64)
	if word >= len(s.words){
		return
	}
	s.words[word] &= ^(1 << bit)
}

func (s *IntSet) Clear(){
	s.words = nil
}

func (s *IntSet) Copy() *IntSet{
	var t IntSet
	for _, word := range s.words {
		t.words = append(t.words, word)
	}
	return &t
}

func (s *IntSet) AddAll(args...int) {
	for _, arg := range args {
		s.Add(arg)
	}
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word & (1 << uint(j)) != 0 {
				if buf.Len() > 1 {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", i * 64 + j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

