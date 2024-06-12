# Recovergoroutine

## Summary
When operating a Go server, we typically use recover middleware. However, if using goroutines, a failure to handle a panic in a goroutine can result in server downtime. To address this, we can use recovergoroutine and check it through linting.

## Supported Go Versions

### 0.2.0 and Below
- Compatible with Go 1.19 and earlier versions.

### 0.3.0 and Above
- Requires Go 1.20 and later versions.

## Install
```bash
go install github.com/Buzzvil/recovergoroutine
```

## Usage
```bash
recovergoroutine -recover="" ./...

# -recover string
#         Custom recovery method name. You can use this option
#         when you want to call a method defined in a struct or
#         use CustomRecover declared in an external package.
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
