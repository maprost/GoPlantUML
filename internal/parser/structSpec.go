package parser

type TypeSpecMap map[string]TypeSpec

type TypeSpec struct {
	Name               string
	Interface          bool
	Relations          []Relation
	FunctionSignatures []string
}

type Relation struct {
	Name        string
	Type        string // primitive or TypeSpec
	Many        bool
	Inheritance bool
}
