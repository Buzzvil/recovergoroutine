package custom

import alias "succdata"

func run() {
	go func() { defer alias.Recover() }()
	go func() { defer handler{}.Recover() }()
}

type handler struct{}

func (handler) Recover() {}

func argumentHandler() {
	go func() { defer consume(alias.Recover) }() // want "recover"
}

func consume(any) {}
