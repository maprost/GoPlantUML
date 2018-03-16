package parser

type TypeSpecMap map[string]*TypeSpec

type TypeSpec struct {
	Name      string
	Interface bool
	Relations []Relation
	Functions []string
}

type Relation struct {
	Name        string
	Type        string // primitive or TypeSpec.Name
	Many        bool
	Inheritance bool
}

type PackageSpec struct {
	Name  string
	Types TypeSpecMap
}

func NewSpec() *PackageSpec {
	return &PackageSpec{
		Types: make(TypeSpecMap),
	}
}

func (p *PackageSpec) AddType(t TypeSpec) {
	if ts, alreadyIn := p.Types[t.Name]; alreadyIn {
		ts.Relations = t.Relations
		ts.Interface = t.Interface

	} else {
		p.Types[t.Name] = &t
	}
}

func (p *PackageSpec) AddFunction(name string, funcSig string) {
	if ts, alreadyIn := p.Types[name]; alreadyIn {
		ts.Functions = append(ts.Functions, funcSig)

	} else {
		p.Types[name] = &TypeSpec{
			Name:      name,
			Functions: []string{funcSig},
		}
	}

}
