package parser

import "sort"

type functionMap map[string][]string

func connectInterfaces(spec *PackageSpec) error {
	structs := createFunctionMap(spec, false)
	interfaces := createFunctionMap(spec, true)

	for strucName, structFuncs := range structs {
		for interfaceName, interfaceFunc := range interfaces {

			i := 0
			for s := 0; s < len(structFuncs) && i < len(interfaceFunc); s++ {
				if structFuncs[s] == interfaceFunc[i] {
					i++
				}
			}

			// all methods in interface are found in the struct
			if i == len(interfaceFunc) {
				spec.Types[strucName].Relations = append(spec.Types[strucName].Relations, Relation{
					Type:        interfaceName,
					Inheritance: true,
				})
			}
		}
	}

	return nil
}

func createFunctionMap(spec *PackageSpec, interfaceType bool) functionMap {
	funcMap := make(functionMap)

	for _, typeSpec := range spec.Types {
		// which type?
		if typeSpec.Interface == interfaceType {

			// has some methods
			if len(typeSpec.Functions) > 0 {
				// copy and sort the functions --> easier to compare
				funcs := make([]string, 0, len(typeSpec.Functions))
				funcs = append(funcs, typeSpec.Functions...)
				sort.Strings(funcs)

				funcMap[typeSpec.Name] = funcs
			}
		}
	}

	return funcMap
}
