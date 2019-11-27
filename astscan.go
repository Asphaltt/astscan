// astscan 测试

package astscan

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/fatih/astrewrite"
)

var (
	_callback Callback
	_contain  Container
)

func getPos(fset *token.FileSet, pos token.Pos) string {
	if fset == nil {
		return ""
	}
	return fset.Position(pos).String()
}

func scan(n ast.Node, fset *token.FileSet, pkg string) (ast.Node, bool) {
	if n == nil {
		return nil, true
	}

	switch v := n.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			if _contain(v.Value) {
				_callback(pkg, getPos(fset, v.ValuePos), v.Value, n)
			}
		}
	case *ast.Comment:
		if _contain(v.Text) {
			_callback(pkg, getPos(fset, v.Slash), v.Text, n)
		}
	}
	return n, true
}

func File(file string, contain Container, callback Callback) error {
	if err := checkParams(contain, callback); err != nil {
		return err
	}

	fset := token.NewFileSet()
	fd, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	astrewrite.Walk(fd, func(n ast.Node) (ast.Node, bool) {
		return scan(n, fset, "")
	})
	return nil
}

func Dir(dir string, contain Container, callback Callback) error {
	if err := checkParams(contain, callback); err != nil {
		return err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for pkg, n := range pkgs {
		astrewrite.Walk(n, func(n ast.Node) (ast.Node, bool) {
			return scan(n, fset, pkg)
		})
	}
	return nil
}

func checkParams(contain Container, callback Callback) error {
	if contain == nil {
		return errors.New("contain cannot be nil")
	}
	if callback == nil {
		return errors.New("callback cannot be nil")
	}

	_contain, _callback = contain, callback
	return nil
}
