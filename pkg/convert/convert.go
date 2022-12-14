package convert

import "strconv"

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) UInt64() (uint64, error) {
	v, err := strconv.Atoi(s.String())
	return uint64(v), err
}

func (s StrTo) MustUInt64() uint64 {
	v, _ := s.UInt64()
	return v
}
