# Recovergoroutine

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/4meepo/tagalign?style=flat-square)

## Summary
When operating a Go server, we typically use recover middleware. However, if using goroutines, a failure to handle a panic in a goroutine can result in server downtime. To address this, we can use recovergoroutine and check it through linting.

## Install
```bash
go install github.com/Buzzvil/recovergoroutine
```

## Usage
```bash
recovergoroutine ./...
```

Check out the test cases for validation [examples](./test/src/faildata/failcode.go).
```bash
# Output:
/Users/scott/Workspace/buzzvil/recovergoroutine/test/src/faildata/failcode.go:4:2: goroutine must have recover
/Users/scott/Workspace/buzzvil/recovergoroutine/test/src/faildata/failcode.go:10:2: goroutine must have recover
/Users/scott/Workspace/buzzvil/recovergoroutine/test/src/faildata/failcode.go:20:2: goroutine must have recover
/Users/scott/Workspace/buzzvil/recovergoroutine/test/src/faildata/failcode.go:25:2: goroutine must have recover
/Users/scott/Workspace/buzzvil/recovergoroutine/test/src/faildata/failcode.go:31:2: goroutine must have recover
/Users/scott/Workspace/buzzvil/recovergoroutine/test/src/faildata/failcode.go:32:2: goroutine must have recover
```

## TODO
- integrate with golangci-lint