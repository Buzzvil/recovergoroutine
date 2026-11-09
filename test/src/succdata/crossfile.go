package succdata

func crossFile() {
	go func() { defer handlePanic() }()
	go func() { defer catcher{}.catch() }()
}

type catcher struct{}

func (catcher) catch() { recover() }
