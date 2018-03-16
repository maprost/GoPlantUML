package main

import (
	"github.com/maprost/GoPlantUML/internal/parser"
)

func main() {

	s := parser.NewSpec()
	err := parser.ReadFile("internal/testdata/test1.go", s)
	if err != nil {
		panic(err)
	}
}
