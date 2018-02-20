package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func ReadFile(file string, structMap TypeSpecMap) error {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return err
	}

	for _, f := range f.Decls {
		// check for  general declarations
		gen, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gen.Tok == token.TYPE {
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// collect structs
				if t, ok := ts.Type.(*ast.StructType); ok {
					sspec := convertStructType(ts.Name.Name, t)
					structMap[sspec.Name] = sspec
				}

				// collect interfaces
				if t, ok := ts.Type.(*ast.InterfaceType); ok {
					sspec := convertInterfaceType(ts.Name.Name, t)
					structMap[sspec.Name] = sspec
				}

			}
		}
	}
	return nil

}

func convertStructType(name string, t *ast.StructType) TypeSpec {
	sspec := TypeSpec{
		Name: name,
	}

	for _, field := range t.Fields.List {
		fmt.Println("Struct Field: ", field.Names, field.Type)
		t, many := convertExpr(field.Type)

		if len(field.Names) == 0 {
			// inheritance
			sspec.Relations = append(sspec.Relations, Relation{
				Name:        "",
				Type:        t,
				Many:        false,
				Inheritance: true,
			})
			continue
		}

		for _, name := range field.Names {
			sspec.Relations = append(sspec.Relations, Relation{
				Name:        name.Name,
				Type:        t,
				Many:        many,
				Inheritance: false,
			})
		}
	}

	return sspec
}

func convertInterfaceType(name string, t *ast.InterfaceType) TypeSpec {
	sspec := TypeSpec{
		Name: name,
	}

	for _, field := range t.Methods.List {
		fmt.Println("Interface Field: ", field.Names, field.Type)
		if len(field.Names) > 1 {
			fmt.Println("Warning found method with weird naming: ", field.Names)
			continue
		}

		funcSig := createFunctionSignature(field.Names[0].Name, field.Type)
		sspec.FunctionSignatures = append(sspec.FunctionSignatures, funcSig)
	}

	return sspec
}

func convertExpr(expr ast.Expr) (string, bool) {
	// array/slice
	if t, ok := expr.(*ast.ArrayType); ok {
		fmt.Println("\tArray", t)
		innerType, _ := convertExpr(t.Elt)
		return innerType, true
	}

	// map
	if t, ok := expr.(*ast.MapType); ok {
		fmt.Println("\tMap", t)
		innerType, _ := convertExpr(t.Value)
		return innerType, true
	}

	// pointer
	if t, ok := expr.(*ast.StarExpr); ok {
		fmt.Println("\tPointer", t)
		innerType, many := convertExpr(t.X)
		return innerType, many
	}

	// struct/interface
	if t, ok := expr.(*ast.Ident); ok {
		fmt.Println("\tType", t)
		return t.Name, false
	}

	// function
	if t, ok := expr.(*ast.FuncType); ok {
		return createFunctionSignature("", t), false
	}

	fmt.Println("\t\tNone of them: ", expr)
	return "", false
}

func createFunctionSignature(name string, expr ast.Expr) string {
	if len(name) > 0 {
		name = " " + name
	}

	// function
	if t, ok := expr.(*ast.FuncType); ok {
		signature := "func" + name +
			"(" + createFunctionFieldString(t.Params.List) + ")" +
			"(" + createFunctionFieldString(t.Results.List) + ")"

		fmt.Println("\tFunction", signature)
		return signature
	}

	return ""
}

func createFunctionFieldString(fields []*ast.Field) string {
	fieldString := ""
	for _, field := range fields {
		fieldType := exprToString(field.Type)
		if len(fieldString) > 0 {
			fieldString += " "
		}
		fieldString += fieldType
	}
	return fieldString
}

func exprToString(expr ast.Expr) string {
	// array/slice
	if t, ok := expr.(*ast.ArrayType); ok {
		str := exprToString(t.Elt)
		return "[]" + str
	}

	// map
	if t, ok := expr.(*ast.MapType); ok {
		keyStr := exprToString(t.Key)
		valStr := exprToString(t.Value)
		return "map[" + keyStr + "]" + valStr
	}

	// pointer
	if t, ok := expr.(*ast.StarExpr); ok {
		str := exprToString(t.X)
		return "*" + str
	}

	// struct/interface
	if t, ok := expr.(*ast.Ident); ok {
		return t.Name
	}

	// function
	if t, ok := expr.(*ast.FuncType); ok {
		return createFunctionSignature("", t)
	}

	return ""
}
