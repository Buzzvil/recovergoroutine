# Recovergoroutine

## Summary
When operating a Go server, we typically use recover middleware. However, if using goroutines, a failure to handle a panic in a goroutine can result in server downtime. To address this, we can use recovergoroutine and check it through linting.

## Supported Go Versions

| Release | Minimum Go |
| --- | --- |
| v0.3.0 | 1.19 |
| v0.4.0 | 1.25 |
| v0.5.0 / main | 1.27.1 |

`-recover`는 신뢰하는 복구 함수 이름입니다. 함수 구현은 검증하지 않습니다.
조건부·늦은 defer 등록의 모든 실행 경로를 검증하지는 않습니다.

## Install
```bash
go install github.com/Buzzvil/recovergoroutine@latest
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
