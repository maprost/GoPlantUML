package parser

import (
	"fmt"
	"go/parser"
	"go/token"
	//"github.com/maprost/toolbox/mpprint"
	"go/ast"
)

func ReadFile(file string) {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Imports:")
	for _, i := range f.Imports {
		fmt.Println(i.Path.Value)
	}

	fmt.Println()
	fmt.Println("Comments:")
	for _, c := range f.Comments {
		fmt.Print(c.Text())
	}

	fmt.Println()
	fmt.Println("Functions:")
	for _, f := range f.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		fmt.Println(fn.Name.Name)
	}

	fmt.Println()
	fmt.Println("Types:")
	for _, f := range f.Decls {
		gen, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		fmt.Println(gen.Tok)
		if gen.Tok == token.TYPE {
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				fmt.Println("TypeSpec: ", ts)
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				fmt.Println("StructType: ", st)
				for _, field := range st.Fields.List {
					fmt.Println("Field: ", field.Names, field.Type)
				}
			}
		}
	}

}
