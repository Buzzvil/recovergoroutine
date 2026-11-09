package custom

import alias "succdata"

func run() {
	go func() { defer alias.Recover() }()
	go func() { defer handler{}.Recover() }()
}

type handler struct{}

func (handler) Recover() {}
