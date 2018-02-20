package parser_test

import (
	"github.com/maprost/GoPlantUML/internal/parser"
	"github.com/maprost/assertion"
	"testing"
)

func TestReadFile(t *testing.T) {
	assert := assertion.New(t)

	spec := make(parser.TypeSpecMap)
	err := parser.ReadFile("../testdata/test1.go", spec)
	assert.Nil(err)

	assert.Contains(spec, parser.TypeSpec{Name: "SubTest"})
	assert.Contains(spec, parser.TypeSpec{Name: "ToInheriate"})
	assert.Contains(spec, parser.TypeSpec{Name: "Test",
		Relations: []parser.Relation{
			{Type: "ToInheriate", Inheritance: true},
			{Name: "counter", Type: "int"},
			{Name: "sub", Type: "SubTest"},
			{Name: "subList", Type: "SubTest", Many: true},
			{Name: "subPointer", Type: "SubTest"},
			{Name: "subPointerList", Type: "SubTest", Many: true},
			{Name: "subListPointer", Type: "SubTest", Many: true},
			{Name: "subMap", Type: "SubTest", Many: true},
			{Name: "subInterface", Type: "SubInterface"},
			{Name: "subFunc", Type: "func(*int)([]*string)"},
			{Name: "subFuncList", Type: "func(int)(string)", Many: true},
		}})
}
