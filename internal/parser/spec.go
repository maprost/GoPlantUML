package parser

// TypeSpecMap definition: struct/interface name --> spec
type TypeSpecMap map[string]*TypeSpec

// TypeSpec holds the information of a struct/interface
type TypeSpec struct {
	Name      string
	Interface bool
	Relations []Relation
	Functions []string
}

// Relation holds every type the struct is using they use
type Relation struct {
	Name        string
	Type        string // primitive or TypeSpec.Name
	Many        bool
	Inheritance bool
}

// PackageSpec holds every information of a package
type PackageSpec struct {
	Name  string
	Types TypeSpecMap
}

// NewSpec returns a new PackageSpec
func NewSpec() *PackageSpec {
	return &PackageSpec{
		Types: make(TypeSpecMap),
	}
}

// AddType add a new TypeSpec, check if there are already function information
func (p *PackageSpec) AddType(t TypeSpec) {
	if ts, alreadyIn := p.Types[t.Name]; alreadyIn {
		ts.Relations = t.Relations
		ts.Interface = t.Interface

	} else {
		p.Types[t.Name] = &t
	}
}

// AddFunction add a function signature, check if there are already type spec information
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
