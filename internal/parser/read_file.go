package parser

import (
	"go/parser"
	"go/token"
)

func ReadFile(file string, spec *PackageSpec) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return err
	}

	err = addFileSpec(f, spec)
	if err != nil {
		return err
	}

	// TODO: load all needed interface specs out of imports (nothing more is needed!)

	return connectInterfaces(spec)
}
