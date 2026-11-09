package recovergoroutine_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Buzzvil/recovergoroutine/recovergoroutine"
)

func TestLint(t *testing.T) {
	t.Run("goroutine has recover", func(t *testing.T) {
		results := analysistest.Run(
			t,
			testDataDir(),
			recovergoroutine.NewAnalyzer(),
			"succdata",
		)
		for _, result := range results {
			assert.NoError(t, result.Err)
		}
	})

	t.Run("goroutine has recover with fail data", func(t *testing.T) {
		results := analysistest.Run(
			t,
			testDataDir(),
			recovergoroutine.NewAnalyzer(),
			"faildata",
		)
		for _, result := range results {
			assert.NoError(t, result.Err)
		}
	})
}

func testDataDir() string {
	_, testFilename, _, ok := runtime.Caller(1)
	if !ok {
		panic("unable to get current test filename")
	}
	return filepath.Join(filepath.Dir(testFilename), "..", "test")
}

func TestCustomRecover(t *testing.T) {
	a := recovergoroutine.NewAnalyzer()
	if err := a.Flags.Set("recover", "Recover"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, testDataDir(), a, "custom")
}
