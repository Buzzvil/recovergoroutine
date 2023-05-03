package faildata

func whenASTFuncLit() {
	go func() {
		defer func() {
			// recover comment
		}()
	}()

	go func() {
		// recover variable
		var recover = 1
		foo(recover)
	}()

	func() {
		// not checked
	}()

	go func() {
		// not used defer
		recover()
	}()

	go func() {
		defer customRecover()
	}()
}

func whenASTIndent() {
	go runGoroutine()
	go nestedFunc1()
}

func runGoroutine() {}

func foo(_ int) {}

func nestedFunc1() {
	// must have recover in parent caller
	recover()
	nestedFunc2()
}

func nestedFunc2() {}

func customRecover() {}
