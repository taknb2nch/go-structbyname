package action

import (
	"fmt"
	"reflect"
)

type My1Action struct{}

func (m *My1Action) Index(s string) string {
	t := reflect.ValueOf(*m).Type()
	fmt.Printf("%s %s, %v\n", t.PkgPath(), t.Name(), s)
	return s + s
}

type My2Action struct{}

func (m *My2Action) Index(s string) string {
	t := reflect.ValueOf(*m).Type()
	fmt.Printf("%s %s, %v\n", t.PkgPath(), t.Name(), s)
	return s + s + s
}
