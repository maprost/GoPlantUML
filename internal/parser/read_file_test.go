package parser_test

import (
	"io/ioutil"
	"testing"

	"github.com/maprost/GoPlantUML/internal/parser"
	"github.com/maprost/assertion"
	"github.com/maprost/toolbox/mpprint"
)

func testFrame(t *testing.T, src string, expected parser.TypeSpecMap) {
	assert := assertion.New(t)

	// prepare test file
	filePath := "../testdata/parserTestFile.go"
	src = "package testdata\n\n" + src
	err := ioutil.WriteFile(filePath, []byte(src), 0644)
	assert.Nil(err)

	// parse
	spec := parser.NewSpec()
	err = parser.ReadFile(filePath, spec)
	assert.Nil(err)

	mpprint.Struct(spec.Types)

	// check
	assert.Len(spec.Types, len(expected), "Different TypeSpecMap size.")
	for _, v := range expected {
		assert.Contains(spec.Types, v)
	}
}

//func TestReadFile(t *testing.T) {
//	assert := assertion.New(t)
//
//	s := parser.NewSpec()
//	err := parser.ReadFile("../testdata/test1.go", parser.ParseFile, s)
//	assert.Nil(err)
//
//	typeSpec := s.Types
//	assert.Contains(typeSpec, parser.TypeSpec{Name: "SubTest"})
//	assert.Contains(typeSpec, parser.TypeSpec{Name: "ToInheriate"})
//	assert.Contains(typeSpec, parser.TypeSpec{Name: "Test",
//		Relations: []parser.Relation{
//			{Type: "ToInheriate", Inheritance: true},
//			{Name: "counter", Type: "int"},
//			{Name: "sub", Type: "SubTest"},
//			{Name: "subList", Type: "SubTest", Many: true},
//			{Name: "subPointer", Type: "SubTest"},
//			{Name: "subPointerList", Type: "SubTest", Many: true},
//			{Name: "subListPointer", Type: "SubTest", Many: true},
//			{Name: "subMap", Type: "SubTest", Many: true},
//			{Name: "subInterface", Type: "SubInterface"},
//			{Name: "subFunc", Type: "func(*int)([]*string)"},
//			{Name: "subFuncList", Type: "func(int)(string)", Many: true},
//		}})
//}

func TestReadFile_simpleStruct(t *testing.T) {
	testFrame(t, `	
		type Blob struct {
			counter int64
		}`,
		parser.TypeSpecMap{
			"": &parser.TypeSpec{
				Name: "Blob",
				Relations: []parser.Relation{
					{Name: "counter", Type: "int64"},
				},
			},
		})
}

func TestReadFile_structWithDifferentTypes(t *testing.T) {
	testFrame(t, `	
		type Blob struct {
			counter         int
			intPointer      *int
			intList         []int
			intListPointer  *[]int
			intPointerList  []*int
			intMap          map[int]string
			simpleFunc      func(int) string
			convertFunc     func(*int) []*string
			complexFuncList []func( func(int)(int, error), map[int]string) func(string)(int,error)
		}`,
		parser.TypeSpecMap{
			"": &parser.TypeSpec{
				Name: "Blob",
				Relations: []parser.Relation{
					{Name: "counter", Type: "int"},
					{Name: "intPointer", Type: "int"},
					{Name: "intList", Type: "int", Many: true},
					{Name: "intListPointer", Type: "int", Many: true},
					{Name: "intPointerList", Type: "int", Many: true},
					{Name: "intMap", Type: "string", Many: true},
					{Name: "simpleFunc", Type: "func(int)(string)"},
					{Name: "convertFunc", Type: "func(*int)([]*string)"},
					{Name: "complexFuncList", Type: "func(func(int)(int,error),map[int]string)(func(string)(int,error))", Many: true},
				},
			},
		})
}

func TestReadFile_structWithDirectInheritance(t *testing.T) {
	testFrame(t, `	
		type Inherit struct {
		}
		type Blob struct {
			Inherit
		}`,
		parser.TypeSpecMap{
			"Inherit": &parser.TypeSpec{
				Name: "Inherit",
			},
			"Blob": &parser.TypeSpec{
				Name: "Blob",
				Relations: []parser.Relation{
					{Name: "", Type: "Inherit", Inheritance: true},
				},
			},
		})
}

func TestReadFile_structWithMultiplyDirectInheritance(t *testing.T) {
	testFrame(t, `	
		type Inherit1 struct {
		}
		type Inherit2 struct {
		}
		type Inherit3 struct {
		}
		type Blob struct {
			Inherit1
			Inherit2
			Inherit3
		}`,
		parser.TypeSpecMap{
			"Inherit1": &parser.TypeSpec{
				Name: "Inherit1",
			},
			"Inherit2": &parser.TypeSpec{
				Name: "Inherit2",
			},
			"Inherit3": &parser.TypeSpec{
				Name: "Inherit3",
			},
			"Blob": &parser.TypeSpec{
				Name: "Blob",
				Relations: []parser.Relation{
					{Type: "Inherit1", Inheritance: true},
					{Type: "Inherit2", Inheritance: true},
					{Type: "Inherit3", Inheritance: true},
				},
			},
		})
}

func TestReadFile_structContainsStructAsVariableType(t *testing.T) {
	testFrame(t, `	
		type Drop struct {
		}
		type Blob struct {
			drop Drop
		}`,
		parser.TypeSpecMap{
			"Drop": &parser.TypeSpec{
				Name: "Drop",
			},
			"Blob": &parser.TypeSpec{
				Name: "Blob",
				Relations: []parser.Relation{
					{Name: "drop", Type: "Drop"},
				},
			},
		})
}

func TestReadFile_interfaceWithFunctions(t *testing.T) {
	testFrame(t, `	
		type Drop interface {
			String() string
			Conv(int, int) string
		}`,
		parser.TypeSpecMap{
			"Drop": &parser.TypeSpec{
				Name:      "Drop",
				Interface: true,
				Functions: []string{
					"func String()(string)",
					"func Conv(int,int)(string)",
				},
			},
		})
}

func TestReadFile_structWithFunctions(t *testing.T) {
	testFrame(t, `	
		type Drop struct {
		}
		func (d Drop) String() string {
			return ""
		}
		func (d *Drop) Conv(int, int) string {
			return ""
		}`,
		parser.TypeSpecMap{
			"Drop": &parser.TypeSpec{
				Name: "Drop",
				Functions: []string{
					"func String()(string)",
					"func Conv(int,int)(string)",
				},
			},
		})
}

func TestReadFile_structWithFunctions_functionsDeclaredFirst(t *testing.T) {
	testFrame(t, `
		func (d Drop) String() string {
			return ""
		}
		func (d *Drop) Conv(int, int) string {
			return ""
		}
		type Drop struct {
		}`,
		parser.TypeSpecMap{
			"Drop": &parser.TypeSpec{
				Name: "Drop",
				Functions: []string{
					"func String()(string)",
					"func Conv(int,int)(string)",
				},
			},
		})
}
