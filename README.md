# Recovergoroutine

## Summary
When operating a Go server, we typically use recover middleware. However, if using goroutines, a failure to handle a panic in a goroutine can result in server downtime. To address this, we can use recovergoroutine and check it through linting.