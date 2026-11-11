package main

import (
	"flag"
	"io"
	"os"
	"strings"

	"github.com/Buzzvil/recovergoroutine/recovergoroutine"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/packages"
)

// main uses syntax loading for ordinary lint runs and delegates driver-specific options to singlechecker.
func main() {
	a := recovergoroutine.NewAnalyzer()
	flags := flag.NewFlagSet(a.Name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	recoverFlag := a.Flags.Lookup("recover")
	flags.Var(recoverFlag.Value, recoverFlag.Name, recoverFlag.Usage)
	tests := flags.Bool("test", true, "include test files")
	if flags.Parse(os.Args[1:]) == nil && flags.NArg() > 0 && !strings.HasSuffix(flags.Arg(0), ".cfg") {
		if code, ok := lintSyntax(a, flags.Args(), *tests); ok {
			os.Exit(code)
		}
	}

	// Let singlechecker parse its own flags from their original defaults.
	singlechecker.Main(recovergoroutine.NewAnalyzer())
}

// lintSyntax returns the lint exit code and whether it handled the invocation.
// On a loading or analysis error, singlechecker takes over to report the failure.
func lintSyntax(a *analysis.Analyzer, patterns []string, tests bool) (int, bool) {
	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.LoadFiles | packages.NeedSyntax | packages.NeedImports,
		Tests: tests,
	}, patterns...)
	if err != nil || len(pkgs) == 0 {
		return 0, false
	}
	for pkg := range packages.Postorder(pkgs) {
		if len(pkg.Errors) != 0 {
			return 0, false
		}
	}

	// The rule reads only syntax. Use the driver's text formatter to preserve
	// positions and deduplicate diagnostics from test variants of a package.
	var graph checker.Graph
	code := 0
	for _, pkg := range pkgs {
		action := &checker.Action{Analyzer: a, Package: pkg, IsRoot: true}
		_, err := a.Run(&analysis.Pass{
			Analyzer: a, Fset: pkg.Fset, Files: pkg.Syntax,
			Report: func(d analysis.Diagnostic) { action.Diagnostics = append(action.Diagnostics, d) },
		})
		if err != nil {
			return 0, false
		}
		if len(action.Diagnostics) > 0 {
			code = 3
			graph.Roots = append(graph.Roots, action)
		}
	}
	if err := graph.PrintText(os.Stderr, -1); err != nil {
		return 1, true
	}
	return code, true
}
