package mystructs

import (
	//"fmt"
	"reflect"
)

import (
	pkg1 "github.com/taknb2nch/go-structbyname/sample/action"
	pkg2 "github.com/taknb2nch/go-structbyname/sample/action/subaction"
)

var structs = make(map[string]reflect.Type)

func init() {
	// action
	register(pkg1.My1Action{})
	register(pkg1.My2Action{})
	register(pkg1.My11Action{})
	register(pkg1.My21Action{})
	// action/subaction
	register(pkg2.My1Action{})
	register(pkg2.My2Action{})
}

func register(x interface{}) {
	t := reflect.TypeOf(x)
	n := t.PkgPath() + "." + t.Name()
	//fmt.Printf("Registered > %v\n", n)
	structs[n] = t
}

func New(name string) (interface{}, bool) {
	t, ok := structs[name]
	if !ok {
		return nil, false
	}
	v := reflect.New(t)
	return v.Interface(), true
}
