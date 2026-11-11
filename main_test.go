package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCLI runs the built binary so diagnostics and exit codes are checked as a CI caller sees them.
func TestCLI(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "recovergoroutine")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}

	for name, src := range map[string]string{
		"go.mod": "module example.org/fixture\n\ngo 1.27\n",
		"p.go": `package fixture
func Recover() {}
func helper() { recover() }
func f() {
	go func() { defer Recover() }()
	go func() { defer helper() }()
	go f()
}
`,
		"p_test.go":  "package fixture\nfunc internal() { go internal() }\n",
		"x_test.go":  "package fixture_test\nfunc external() { go external() }\n",
		"ignored.go": "//go:build ignore\n\npackage fixture\nfunc ignored() { go ignored() }\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}

	run := func(t *testing.T, debug bool, wantCode int, args ...string) string {
		t.Helper()
		cmd := exec.Command(bin, args...)
		cmd.Dir = dir
		cmd.Env = append(cmd.Environ(), "GOWORK=off", "GOFLAGS=", "GOPACKAGESDRIVER=off", "GOPACKAGESDEBUG=false")
		if debug {
			cmd.Env = append(cmd.Env, "GOPACKAGESDEBUG=true")
		}
		out, err := cmd.CombinedOutput()
		var exit *exec.ExitError
		if err != nil && !errors.As(err, &exit) {
			t.Fatal(err)
		}
		if cmd.ProcessState.ExitCode() != wantCode {
			t.Fatalf("%v: got exit %d, want %d\n%s", args, cmd.ProcessState.ExitCode(), wantCode, out)
		}
		return string(out)
	}

	t.Run("syntax load", func(t *testing.T) {
		out := run(t, true, 3, "-recover=Recover", "./...")
		if !strings.Contains(out, "-export=false") || strings.Contains(out, "-export=true") {
			t.Fatalf("expected package loading without dependency export data:\n%s", out)
		}
	})
	t.Run("type error without lint diagnostics", func(t *testing.T) {
		fixture := filepath.Join(dir, "testdata")
		if err := os.Mkdir(fixture, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(fixture, "broken.go"), []byte("package broken\nvar _ int = false\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if out := run(t, false, 0, "./testdata"); out != "" {
			t.Fatalf("expected no diagnostics, got %q", out)
		}
	})
	t.Run("text diagnostics", func(t *testing.T) {
		out := run(t, false, 3, "-recover", "Recover", "./...")
		if strings.Count(out, "\n") != 3 {
			t.Fatalf("expected exactly three diagnostics:\n%s", out)
		}
		for _, pos := range []string{"p.go:7:2", "p_test.go:2:19", "x_test.go:2:19"} {
			if !strings.Contains(out, pos+": use function literals when using goroutines") {
				t.Errorf("missing %s in:\n%s", pos, out)
			}
		}
	})
	t.Run("exclude tests", func(t *testing.T) {
		out := run(t, false, 3, "-test=false", ".")
		if strings.Count(out, "\n") != 2 || !strings.Contains(out, "p.go:5:2: goroutine must have recover") || !strings.Contains(out, "p.go:7:2: use function literals when using goroutines") {
			t.Fatal(out)
		}
	})
	t.Run("missing package", func(t *testing.T) {
		out := run(t, false, 1, "./nonexistent")
		if !strings.Contains(out, "nonexistent") {
			t.Fatal(out)
		}
	})

	t.Run("help after recover", func(t *testing.T) {
		out := run(t, false, 0, "-recover=Recover", "-h")
		if !strings.Contains(out, "Usage:") || strings.Contains(out, `(default "Recover")`) {
			t.Fatal(out)
		}
	})
	t.Run("singlechecker flags", func(t *testing.T) {
		out := run(t, false, 0, "-json", "-recover=Recover", ".")
		if !strings.Contains(out, `"recovergoroutine"`) || !strings.Contains(out, "use function literals when using goroutines") {
			t.Fatal(out)
		}
	})
}
