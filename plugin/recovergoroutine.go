package main

import (
	linters "github.com/Buzzvil/recovergoroutine/recovergoroutine"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		linters.NewAnalyzer(),
	}
}

var AnalyzerPlugin analyzerPlugin
