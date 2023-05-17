package succdata

func whenASTFuncLit() {
	go func() {
		defer recover()
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
		rec := func() {
			defer func() {
				recover()
			}()
		}

		defer rec()
	}()
}
