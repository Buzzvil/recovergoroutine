package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Buzzvil/recovergoroutine/recovergoroutine"
)

func main() {
	singlechecker.Main(recovergoroutine.Analyzer)
}
