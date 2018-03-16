package parser

import (
	"fmt"
	"go/ast"
)

func addFileSpec(file *ast.File, specMap *PackageSpec) error {
	for _, f := range file.Decls {

		// GenDecl contains structs/interfaces/vars/imports (declarations with list behavior)
		if gen, ok := f.(*ast.GenDecl); ok {

			for _, genSpec := range gen.Specs {

				// TypeSpec contains structs/interface
				if ts, ok := genSpec.(*ast.TypeSpec); ok {

					// collect structs
					if t, ok := ts.Type.(*ast.StructType); ok {
						typeSpec := convertStructType(ts.Name.Name, t)
						specMap.AddType(typeSpec)
					}

					// collect interfaces
					if t, ok := ts.Type.(*ast.InterfaceType); ok {
						typeSpec := convertInterfaceType(ts.Name.Name, t)
						specMap.AddType(typeSpec)
					}
				}
			}
		}

		// FuncDecl contains functions
		if fun, ok := f.(*ast.FuncDecl); ok {
			// function has a receiver -> belongs to struct
			if len(fun.Recv.List) > 0 {
				structName, _ := convertExpr(fun.Recv.List[0].Type)
				funcSig := createFunctionSignature(fun.Name.Name, fun.Type)

				specMap.AddFunction(structName, funcSig)
			}
		}
	}
	return nil
}

func convertStructType(name string, t *ast.StructType) TypeSpec {
	typeSpec := TypeSpec{
		Name:      name,
		Interface: false,
	}

	for _, field := range t.Fields.List {
		fmt.Println("Struct Field: ", field.Names, field.Type)
		t, many := convertExpr(field.Type)

		if len(field.Names) == 0 {
			// inheritance
			typeSpec.Relations = append(typeSpec.Relations, Relation{
				Name:        "",
				Type:        t,
				Many:        false,
				Inheritance: true,
			})
			continue
		}

		for _, name := range field.Names {
			typeSpec.Relations = append(typeSpec.Relations, Relation{
				Name:        name.Name,
				Type:        t,
				Many:        many,
				Inheritance: false,
			})
		}
	}

	return typeSpec
}

func convertInterfaceType(name string, t *ast.InterfaceType) TypeSpec {
	typeSpec := TypeSpec{
		Name:      name,
		Interface: true,
	}

	for _, field := range t.Methods.List {
		fmt.Println("Interface Field: ", field.Names, field.Type)
		if len(field.Names) > 1 {
			fmt.Println("Warning found method with weird naming: ", field.Names)
			continue
		}

		funcSig := createFunctionSignature(field.Names[0].Name, field.Type)
		typeSpec.Functions = append(typeSpec.Functions, funcSig)
	}

	return typeSpec
}

func convertExpr(expr ast.Expr) (innerType string, many bool) {
	// array/slice
	if t, ok := expr.(*ast.ArrayType); ok {
		fmt.Println("\tArray", t)
		innerType, _ = convertExpr(t.Elt)
		return innerType, true
	}

	// map
	if t, ok := expr.(*ast.MapType); ok {
		fmt.Println("\tMap", t)
		innerType, _ = convertExpr(t.Value)
		return innerType, true
	}

	// pointer
	if t, ok := expr.(*ast.StarExpr); ok {
		fmt.Println("\tPointer", t)
		innerType, many = convertExpr(t.X)
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
			fieldString += ","
		}
		fieldString += fieldType
	}
	return fieldString
}
