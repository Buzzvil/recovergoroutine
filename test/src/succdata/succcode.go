package succdata

func whenASTFuncLit() {
	go func() {
		defer func() { recover() }()
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
			}
		}()
	}()

	go func() {
		defer func() {
			recover()
		}()
	}()

	go func() {
		defer handlePanic()
	}()
}

func handlePanic() { recover() }

func Recover() { recover() }
