package recovergoroutine

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer() *analysis.Analyzer {
	analyzer := &analysis.Analyzer{
		Name: "recovergoroutine",
		Doc:  "finds goroutine code without recover",
		Run:  run,
	}

	return analyzer
}

func run(pass *analysis.Pass) (interface{}, error) {
	var runErr error
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			goStmt, ok := n.(*ast.GoStmt)
			if !ok {
				return true
			}

			ok, err := safeGoStmt(goStmt, pass)
			if err != nil {
				runErr = err
				return false
			}

			if ok {
				return true
			}

			pass.Report(analysis.Diagnostic{
				Pos:      goStmt.Pos(),
				End:      0,
				Category: "goroutine",
				Message:  "goroutine must have recover",
			})

			return false
		})
	}

	return nil, runErr
}

func safeGoStmt(goStmt *ast.GoStmt, pass *analysis.Pass) (bool, error) {
	fn := goStmt.Call
	switch fun := fn.Fun.(type) {
	case *ast.SelectorExpr:
		return safeSelectorExpr(fun, pass, safeFunc)
	case *ast.FuncLit:
		return safeFunc(fun, pass)
	case *ast.Ident:
		if fun.Obj == nil {
			return false, nil
		}

		funcDecl, ok := fun.Obj.Decl.(*ast.FuncDecl)
		if !ok {
			return false, nil
		}

		return safeFunc(funcDecl, pass)
	}

	return false, fmt.Errorf("unexpected goroutine function type: %v", reflect.TypeOf(fn.Fun).String())
}

func safeFunc(node ast.Node, pass *analysis.Pass) (bool, error) {
	result := false
	var err error
	ast.Inspect(node, func(node ast.Node) bool {
		deferStmt, ok := node.(*ast.DeferStmt)
		if !ok {
			return true
		}

		ok, err = hasRecover(deferStmt.Call, pass)
		if err != nil {
			return false
		}

		if ok {
			result = true
			return false
		}

		return !result
	})

	return result, err
}

func hasRecover(expr ast.Node, pass *analysis.Pass) (bool, error) {
	var result bool
	var err error
	ast.Inspect(expr, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CallExpr:
			if isRecover(n) || isCustomRecover(n) {
				result = true
				return false
			}
		case *ast.SelectorExpr:
			if n.Sel == nil {
				return true
			}

			var ok bool
			ok, err = safeSelectorExpr(n, pass, hasRecover)
			if err != nil {
				return false
			}

			if ok {
				result = true
				return false
			}
		}
		return true
	})

	return result, err
}

func safeSelectorExpr(
	expr *ast.SelectorExpr,
	pass *analysis.Pass,
	methodChecker func(node ast.Node, pass *analysis.Pass) (bool, error),
) (bool, error) {
	ident, ok := expr.X.(*ast.Ident)
	if !ok {
		return false, nil
	}

	methodName := expr.Sel.Name
	objType := pass.TypesInfo.ObjectOf(ident)
	pointerType, ok := objType.Type().(*types.Pointer)
	if !ok {
		return false, nil
	}

	named, ok := pointerType.Elem().(*types.Named)
	if !ok {
		return false, nil
	}

	result := false
	for i := 0; i < named.NumMethods(); i++ {
		if named.Method(i).Name() != methodName {
			continue
		}

		fset := pass.Fset
		position := fset.Position(named.Method(i).Pos())
		file, err := parser.ParseFile(fset, position.Filename, nil, 0)
		if err != nil {
			return false, fmt.Errorf("parse file: %w", err)
		}

		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == methodName {
					result, err = methodChecker(funcDecl, pass)
					break
				}
			}
		}
	}

	return result, nil
}

func isRecover(callExpr *ast.CallExpr) bool {
	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "recover"
}

func isCustomRecover(callExpr *ast.CallExpr) bool {
	result := false
	switch fun := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		if fun.Sel == nil {
			return result
		}

		result = checkIdent(fun.Sel)
	case *ast.Ident:
		result = checkIdent(fun)
	}

	return result
}

func checkIdent(ident *ast.Ident) bool {
	result := false
	if ident.Obj == nil {
		return result
	}

	funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl)
	if !ok {
		return result
	}

	ast.Inspect(funcDecl, func(node ast.Node) bool {
		if callExpr, ok := node.(*ast.CallExpr); ok && isRecover(callExpr) {
			result = true
			return false
		}

		return true
	})

	return result
}
