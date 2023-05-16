package recovergoroutine_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Buzzvil/recovergoroutine/recovergoroutine"
)

type CustomTesting struct {
	t *testing.T
}

func (c *CustomTesting) Errorf(format string, args ...any) {
	c.t.Log(fmt.Sprintf(format, args...))
}

func TestLint(t *testing.T) {
	t.Run("goroutine has recover", func(t *testing.T) {
		results := analysistest.Run(
			&CustomTesting{t: t},
			testDataDir(),
			recovergoroutine.NewAnalyzer(),
			"succdata",
		)
		for _, result := range results {
			assert.Len(t, result.Diagnostics, 0)
			assert.NoError(t, result.Err)
		}
	})

	t.Run("goroutine has recover with fail data", func(t *testing.T) {
		results := analysistest.Run(
			&CustomTesting{t: t},
			testDataDir(),
			recovergoroutine.NewAnalyzer(),
			"faildata",
		)
		for _, result := range results {
			assert.Len(t, result.Diagnostics, 8)
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
