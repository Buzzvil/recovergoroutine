package faildata

func whenASTFuncLit() {
	go func() { // want "recover"
		defer func() {
			// recover comment
		}()
	}()

	go func() { // want "recover"
		// recover variable
		var recover = 1
		foo(recover)
	}()

	func() {
		// not checked
	}()

	go func() { // want "recover"
		// not used defer
		recover()
	}()

	go func() { // want "recover"
		defer customRecover()
	}()
}

func whenASTIndent() {
	go runGoroutine() // want "function literals"
	go nestedFunc1()  // want "function literals"
}

func whenCallMethod() {
	foo := &Foo{}
	go foo.run() // want "function literals"
	go func() {  // want "recover"
		defer foo.Recover()
	}()
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

type Foo struct{}

func (a *Foo) run() {}

func (a *Foo) Recover() {}

func nestedHandlers() {
	go func() { // want "recover"
		go func() { defer func() { recover() }() }()
	}()
}

func nestedMissing() {
	go func() { // want "recover"
		go func() {}() // want "recover"
	}()
}

func bareRecover() {
	go func() { defer recover() }() // want "recover"
}

func argumentRecover() {
	go func() { defer consume(recover()) }() // want "recover"
}

func consume(any) {}

func shadowedRecover() {
	recover := func() {}
	go func() { defer func() { recover() }() }() // want "recover"
}

func indirectRecover() {
	go func() { // want "recover"
		defer func() { nestedFunc1() }()
	}()
	go func() { // want "recover"
		defer func() { func() { recover() }() }()
	}()
}
