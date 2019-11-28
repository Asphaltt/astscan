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
	_check    Checker
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
			if _check(v.Value) {
				_callback(Item{pkg, getPos(fset, v.ValuePos), TypeString, v.Value, n})
			}
		}
	case *ast.Comment:
		if _check(v.Text) {
			_callback(Item{pkg, getPos(fset, v.Slash), TypeComment, v.Text, n})
		}
	}
	return n, true
}

func File(file string, contain Checker, callback Callback) error {
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

func Dir(dir string, contain Checker, callback Callback) error {
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

func checkParams(contain Checker, callback Callback) error {
	if contain == nil {
		return errors.New("contain cannot be nil")
	}
	if callback == nil {
		return errors.New("callback cannot be nil")
	}

	_check, _callback = contain, callback
	return nil
}
