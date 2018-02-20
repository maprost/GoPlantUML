package testdata

import "strconv"

type SubInterface interface {
	String() string
	Blob(func(*string) []*int, int) map[*int]**SubInterface
}

type SubTest struct {
}

type ToInheriate struct {
}

type Test struct {
	ToInheriate
	counter        int
	sub            SubTest
	subList        []SubTest
	subPointer     *SubTest
	subPointerList *[]SubTest
	subListPointer []*SubTest
	subMap         map[string]SubTest
	subInterface   SubInterface
	subFunc        func(*int) []*string
	subFuncList    []func(int) string
}

func (t *Test) String() string {
	return strconv.Itoa(t.counter)
}
