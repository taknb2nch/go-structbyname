package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"go/ast"
	"go/parser"
	"go/token"

	//myruntime "./runtime"
)

const (
	OWN_NAME = "generator.go"

	NEWLINE_FLAG_U = "u"
	NEWLINE_FLAG_W = "w"
)

var (
	paramRootDir       = flag.String("pd", "", "project root directory : required")
	paramPackageRoot   = flag.String("pr", "", "package root (import path) : optional")
	paramOutput        = flag.String("o", "./mystructs/structs.go", "output file (from project root) : optional")
	paramStructDir     = flag.String("d", "./", "target directory (from project root) : optional")
	paramLineSeparetor = flag.String("L", "", "line separetor u: LF, w: CRLF, default: auto detect : optional")
)

func main() {
	var (
		rootDir     string
		packageRoot string
		outputFile  string
		structDir   string
	)
	var err error

	flag.Parse()

	if *paramRootDir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	rootDir, err = filepath.Abs(*paramRootDir)
	if err != nil {
		panic(err)
	}

	packageRoot = *paramPackageRoot
	if packageRoot == "" {
		packageRoot, err = getProjectRoot(rootDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
	if !strings.HasSuffix(packageRoot, "/") {
		packageRoot = packageRoot + "/"
	}

	outputFile = filepath.Join(rootDir, *paramOutput)

	_, err = os.Stat(outputFile)
	if err == nil {
		if err := os.RemoveAll(outputFile); err != nil {
			panic(err)
		}
		fmt.Println("deleted", outputFile)
	}

	structDir = filepath.Join(rootDir, *paramStructDir)

	if *paramLineSeparetor != "" &&
		*paramLineSeparetor != NEWLINE_FLAG_U &&
		*paramLineSeparetor != NEWLINE_FLAG_W {

		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("packageRoot", packageRoot)
	fmt.Println("outputFile", outputFile)
	fmt.Println("structDir", structDir)

	mm := make(map[string][]string)
	fset := token.NewFileSet()
	filepath.Walk(structDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		if info.Name() == OWN_NAME {
			return nil
		}

		p := replacePath(rootDir, path)

		pkgs, err := parser.ParseDir(fset, path, nil, 0)

		if err != nil {
			panic(err)
		}

		for _, v := range pkgs {
			for _, file := range v.Files {
				ss := parseFile(file)

				if len(ss) == 0 {
					continue
				}

				ms := mm[p]
				ms = append(ms, ss...)
				mm[p] = ms
			}
		}

		return nil
	})

	i := 1
	keys := sortedKeys(mm)

	mds := make([]Pkg, 0)

	for _, key := range keys {
		md := Pkg{
			fmt.Sprintf("pkg%d", i),
			key,
			mm[key],
		}

		fmt.Println(key, mm[key])

		mds = append(mds, md)
		i++
	}

	output(outputFile,
		MyData{
			MyPackage:   parseMyPackage(*paramOutput),
			BasePackage: packageRoot,
			Packages:    mds,
		},
		*paramLineSeparetor)
}

func getProjectRoot(rootDir string) (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", errors.New("GOPATH not found.")
	}

	gopath = filepath.Join(gopath, "src")

	if !strings.HasPrefix(rootDir, gopath) {
		return "", errors.New("Can not find project root.")
	}

	return filepath.ToSlash(rootDir[len(gopath)+1:]), nil
}

func sortedKeys(m map[string][]string) []string {
	mk := make([]string, len(m))
	i := 0
	for k, _ := range m {
		mk[i] = k
		i++
	}
	sort.Strings(mk)

	return mk
}

type MyData struct {
	MyPackage   string
	BasePackage string
	Packages    []Pkg
}

type Pkg struct {
	Alias      string
	ImportPath string
	Structs    []string
}

func mkdir(path string) string {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); err != nil {
		os.MkdirAll(dir, 0777)
	}

	return dir
}

func parseMyPackage(path string) string {
	dir := filepath.Dir(path)

	if dir == "." {
		return "main"
	}

	ps := strings.Split(filepath.ToSlash(dir), "/")

	return ps[len(ps)-1]
}

func output(outpath string, mds MyData, newlineFlag string) {
	var err error

	mkdir(outpath)

	tpl := template.Must(template.New("mytemplate").Parse(template_text))

	b := make([]byte, 0)
	buf := bytes.NewBuffer(b)

	if err = tpl.Execute(buf, mds); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	s := strings.TrimSpace(buf.String())
	s, err = convertLineSeparetor(s, newlineFlagToEnum(newlineFlag))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(outpath, []byte(s), 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// func output(outpath string, mds MyData, newlineFlag string) {
// 	mkdir(outpath)

// 	file, err := os.OpenFile(outpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n", err)
// 		os.Exit(1)
// 	}

// 	tpl := template.Must(template.New("mytemplate").Parse(template_text))

// 	if err := tpl.Execute(myruntime.NewWriter(file, newlineFlagToEnum(newlineFlag)), mds); err != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n", err)
// 		file.Close()
// 		os.Exit(1)
// 	}

// 	if err = file.Close(); err != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n", err)
// 		os.Exit(1)
// 	}
// }

func parseFile(f *ast.File) []string {
	ss := make([]string, 0)
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range gd.Specs {
			ts, ok1 := spec.(*ast.TypeSpec)
			if !ok1 {
				continue
			}

			_, ok2 := ts.Type.(*ast.StructType)
			if !ok2 {
				continue
			}

			ss = append(ss, ts.Name.Name)
		}
	}

	return ss
}

func replacePath(rootDir, path string) string {
	if !strings.HasPrefix(path, rootDir) {
		panic(fmt.Sprintf("invalid path: %s", path))
	}

	if rootDir == path {
		return "./"
	} else {
		return filepath.ToSlash(path[len(rootDir)+1:])
	}
}

const template_text = `package {{$.MyPackage}}

import (
	//"fmt"
	"reflect"
)

import (
	{{range $i, $md := .Packages}}{{if $i}}
	{{end}}{{$md.Alias}} "{{$.BasePackage}}{{$md.ImportPath}}"{{end}}
)

var structs = make(map[string]reflect.Type)

func init() {
	{{range $i, $md := .Packages}}{{if $i}}
	{{end}}// {{$md.ImportPath}}
	{{range $ii, $st := $md.Structs}}{{if $ii}}
	{{end}}register({{$md.Alias}}.{{$st}}{}){{end}}{{end}}
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
`

type PlatformType uint32

const (
	AutoDetect PlatformType = iota
	Unix
	Windows
)

func convertLineSeparetor(s string, platformType PlatformType) (string, error) {
	if platformType != Unix && platformType != Windows {
		if runtime.GOOS == "windows" {
			platformType = Windows
		} else {
			platformType = Unix
		}
	}

	var (
		r   *regexp.Regexp
		err error
	)

	if platformType == Windows {
		r, err = regexp.Compile("\r\n")
		if err != nil {
			return s, err
		}

		s = r.ReplaceAllString(s, "\n")

		r, err = regexp.Compile("\n")
		if err != nil {
			return s, err
		}

		return r.ReplaceAllString(s, "\r\n"), nil
	} else {
		r, err = regexp.Compile("\r\n")
		if err != nil {
			return s, err
		}

		return r.ReplaceAllString(s, "\n"), nil
	}
}

func newlineFlagToEnum(s string) PlatformType {
	switch s {
	case NEWLINE_FLAG_U:
		return Unix
	case NEWLINE_FLAG_W:
		return Windows
	default:
		return AutoDetect
	}
}

// func newlineFlagToEnum(s string) myruntime.PlatformType {
// 	switch s {
// 	case NEWLINE_FLAG_U:
// 		return myruntime.Unix
// 	case NEWLINE_FLAG_W:
// 		return myruntime.Windows
// 	default:
// 		return myruntime.AutoDetect
// 	}
// }
