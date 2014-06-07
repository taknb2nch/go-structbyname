package main

import (
	"errors"
	"fmt"
	"reflect"

	"./mystructs"
)

func main() {
	execute("github.com/taknb2nch/go-structbyname/sample/action.My1Action", "Index")
	execute("github.com/taknb2nch/go-structbyname/sample/action/subaction.My1Action", "Index")

	execute("github.com/taknb2nch/go-structbyname/sample/action.My1000Action", "Index")
	execute("github.com/taknb2nch/go-structbyname/sample/action.My1Action", "Index2")
}

func execute(structName, methodName string) error {
	a1, ok := mystructs.New(structName)
	if !ok {
		return errors.New(structName + " is not found.")
	}

	//fmt.Println(a1, ok)

	v := reflect.ValueOf(a1)
	m := v.MethodByName(methodName)

	//fmt.Println(v, m)

	if !m.IsValid() || m.IsNil() {
		return errors.New(structName + "." + methodName + " is not found.")
	}

	vs := m.Call([]reflect.Value{
		reflect.ValueOf("aaa"),
	})

	for _, vss := range vs {
		fmt.Println(vss)
	}

	return nil
}
