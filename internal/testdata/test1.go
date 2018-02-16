package testdata

import "strconv"

type Test struct {
	counter int
}

func (t *Test) String () string {
	return strconv.Itoa(t.counter)
}