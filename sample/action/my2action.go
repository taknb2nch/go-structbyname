package action

import (
	"fmt"
	"reflect"
)

type My11Action struct{}

func (m *My11Action) Index(s string) string {
	t := reflect.ValueOf(*m).Type()
	fmt.Printf("%s %s, %v\n", t.PkgPath(), t.Name(), s)
	return s + s + s + s
}

type My21Action struct{}

func (m *My21Action) Index(s string) string {
	t := reflect.ValueOf(*m).Type()
	fmt.Printf("%s %s, %v\n", t.PkgPath(), t.Name(), s)
	return s + s + s + s + s
}
