package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type visitor struct {
	imports    []*ast.ImportSpec
	dropVars   []*ast.ValueSpec
	interfaces []*ast.InterfaceType
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			if s, ok := spec.(*ast.ImportSpec); ok {
				v.imports = append(v.imports, s)
			} else if s, ok := spec.(*ast.ValueSpec); ok {
				for _, ident := range s.Names {
					if ident.Name == "_" { // 仅取 var _ pacakge.Type
						v.dropVars = append(v.dropVars, s)
					}
				}
			}
		}
	case *ast.InterfaceType:
		v.interfaces = append(v.interfaces, n)
	}

	return v
}

type Parser struct {
	v *visitor
}

func (p *Parser) Parse(srcFile string) error {
	var err error
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	p.v = &visitor{}
	ast.Walk(p.v, f)
	return err
}

func (p *Parser) GetImportSpec() []*ast.ImportSpec {
	return p.v.imports
}

func (p *Parser) GetDropVarSpec() []*ast.ValueSpec {
	return p.v.dropVars
}

func (p *Parser) GetInterfaces() []*ast.InterfaceType {
	return p.v.interfaces
}

var srcFile = flag.String("f", "", "help message for flagname")

func main() {
	flag.Parse()
	if *srcFile == "" {
		fmt.Println("please add src file use -f")
		return
	}
	p := Parser{}
	p.Parse(*srcFile)
	vars := p.GetDropVarSpec()
	ast.Print(nil, vars)
	imports := p.GetImportSpec()
	ast.Print(nil, imports)
	interfaces := p.GetInterfaces()
	ast.Print(nil, interfaces)
}
