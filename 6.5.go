package main

import (
	"fmt"
	"bytes"
)

const platForm = 32 << (^uint(0) >> 63)

func main() {
	fmt.Printf("%T\n", platForm)
	var s IntSet
	s.AddAll(10, 20, 30)
	fmt.Printf("%v\n", s.Elems())
}

type IntSet struct {
	words []uint
}

func (s *IntSet) Has(x int) bool {
	word, bit := x / platForm, uint(x % platForm)
	return word < len(s.words) && s.words[word] & (1 << bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x / platForm, uint(x % platForm)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

//并集
func (s *IntSet) UnionWith(t *IntSet) {
	for i, word := range t.words {
		if i < len(s.words) {
			s.words[i] |= word	
		}else {
			s.words = append(s.words, word)
		}
	}
}
//交集
func (s *IntSet) IntersectWith(t *IntSet) *IntSet {
	var r IntSet
	for i, word := range s.words {
		if i < len(t.words) {
			r.words = append(r.words, word & t.words[i])
		}
	}
	return &r
}
//差集
func (s *IntSet) DifferenceWith(t *IntSet) *IntSet {
	var r IntSet
	for i, word := range s.words {
		if i < len(t.words) {
			r.words = append(r.words, word & ^t.words[i])
		} else{
			r.words = append(r.words, word)
		}
	}

	return &r
}
//并差集
func (s *IntSet) SymmetricDifference(t *IntSet) *IntSet {
	var r IntSet

	if len(s.words) < len(t.words) {
		s, t = t, s
	}
	for i, word := range s.words {
		if i < len(t.words) {
			r.words = append(r.words, word ^ t.words[i])
		} else{
			r.words = append(r.words, word)
		}
	}

	return &r
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
	word, bit := x / platForm, uint(x % platForm)
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
		for j := 0; j < platForm; j++ {
			if word & (1 << uint(j)) != 0 {
				if buf.Len() > 1 {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", i * platForm + j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func (s *IntSet) Elems() []int {
	var r []int
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < platForm; j++ {
			if word & (1 << uint(j)) != 0 {
				r = append(r, i * platForm + j)
			}
		}
	}
	return r
}
