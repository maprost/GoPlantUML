package parser_test

import (
	"testing"

	"github.com/maprost/GoPlantUML/internal/parser"
)

func TestReadFile_structInheritInterface(t *testing.T) {
	testFrame(t, `
		type Blob interface {
			String() string
		}
		type Drop struct {
		}
		func (d Drop) String() string {
			return ""
		}`,
		parser.TypeSpecMap{
			"Blob": &parser.TypeSpec{
				Name:      "Blob",
				Interface: true,
				Functions: []string{
					"func String()(string)",
				},
			},
			"Drop": &parser.TypeSpec{
				Name: "Drop",
				Relations: []parser.Relation{
					{Type: "Blob", Inheritance: true},
				},
				Functions: []string{
					"func String()(string)",
				},
			},
		})
}

func TestReadFile_structInheritInterface_multiplyMethods(t *testing.T) {
	testFrame(t, `
		type Blob interface {
			String() string
		}
		type Drop struct {
		}
		func (d Drop) Loop() string {
			return ""
		}
		func (d Drop) String() string {
			return ""
		}
		func (d Drop) Boom() string {
			return ""
		}`,
		parser.TypeSpecMap{
			"Blob": &parser.TypeSpec{
				Name:      "Blob",
				Interface: true,
				Functions: []string{
					"func String()(string)",
				},
			},
			"Drop": &parser.TypeSpec{
				Name: "Drop",
				Relations: []parser.Relation{
					{Type: "Blob", Inheritance: true},
				},
				Functions: []string{
					"func Loop()(string)",
					"func String()(string)",
					"func Boom()(string)",
				},
			},
		})
}

func TestReadFile_structCanNotInheritInterface_doNotHaveAllMethods(t *testing.T) {
	testFrame(t, `
		type Blob interface {
			String() string
			Boom() string
		}
		type Drop struct {
		}
		func (d Drop) String() string {
			return ""
		}`,
		parser.TypeSpecMap{
			"Blob": &parser.TypeSpec{
				Name:      "Blob",
				Interface: true,
				Functions: []string{
					"func String()(string)",
					"func Boom()(string)",
				},
			},
			"Drop": &parser.TypeSpec{
				Name: "Drop",
				Functions: []string{
					"func String()(string)",
				},
			},
		})
}

func TestReadFile_structCanNotInheritInterface_interfaceHasNoMethods(t *testing.T) {
	testFrame(t, `
		type Blob interface {
		}
		type Drop struct {
		}
		func (d Drop) String() string {
			return ""
		}`,
		parser.TypeSpecMap{
			"Blob": &parser.TypeSpec{
				Name:      "Blob",
				Interface: true,
			},
			"Drop": &parser.TypeSpec{
				Name: "Drop",
				Functions: []string{
					"func String()(string)",
				},
			},
		})
}

func TestReadFile_structInheritInterface_multiplyMethodsForBoth(t *testing.T) {
	testFrame(t, `
		type Blob interface {
			String() string
			Loop() string
		}
		type Drop struct {
		}
		func (d Drop) Loop() string {
			return ""
		}
		func (d Drop) String() string {
			return ""
		}
		func (d Drop) Boom() string {
			return ""
		}`,
		parser.TypeSpecMap{
			"Blob": &parser.TypeSpec{
				Name:      "Blob",
				Interface: true,
				Functions: []string{
					"func String()(string)",
					"func Loop()(string)",
				},
			},
			"Drop": &parser.TypeSpec{
				Name: "Drop",
				Relations: []parser.Relation{
					{Type: "Blob", Inheritance: true},
				},
				Functions: []string{
					"func Loop()(string)",
					"func String()(string)",
					"func Boom()(string)",
				},
			},
		})
}
