go-structbyname
===============

source code generator to create a instance of struct by name in golang.

##Usage
```$ go run generate.go -pd [project root directory]```

###Options

- -pd : project root directory (required)
- -pr : package root (optional)
- -o  : output file (optional: ./mystructs/structs.go)
- -d  : target directory to parse (optional: ./)

If you run the generate command, mystructs/structs.go is created. 
You can get an instance of the structure by importing mystructs pacakge, and by executing mystruct.New function with the structure name.


###Simple Implementation
```go
package main

import (
	"fmt"
	"reflect"

	"./mystructs"
)

func main() {
	structName := "github.com/taknb2nch/go-structbyname/sample/action.My1Action"
	methodName := "Index1"

	a1, ok := mystructs.New(structName)
	if !ok {
		fmt.Printf("%s is not found.", structName)
		return
	}

	v := reflect.ValueOf(a1)
	m := v.MethodByName(methodName)

	if !m.IsValid() || m.IsNil() {
		fmt.Printf("%s.%s is not found.", structName, methodName)
		return
	}

	vs := m.Call([]reflect.Value{
		reflect.ValueOf("aaa"),
	})

	for _, vss := range vs {
		fmt.Println(vss)
	}
}
```

##Sample
There is a sample in sample directory. 
mystructs and structs.go are created by generator.  
If you try,

1. delete mystructs directory
1. execute generator.   
```$ go run generate.go -pd ./sample```
1. execute main  
```$ go run ./sample/main.go```

##License
[MIT License](https://github.com/taknb2nch/go-structbyname/blob/master/LICENSE)

