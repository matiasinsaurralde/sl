package main

import (
	"github.com/matiasinsaurralde/sl/runtime"
)

func main() {
	runtime, err := SL.NewRuntime("../test/single_declaration.sl")
	if err != nil {
		panic(err)
	}
	runtime.Start()
}
