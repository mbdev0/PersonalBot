package program

import (
	"sync"
)

var (
	programMu sync.RWMutex
	programs  = make(map[string]Program)
)

// Register makes a Program available by the provided name.
// If Register is called twice with the same name or if program is nil,
// it panics.
func Register(name string, program Program) {
	programMu.Lock()
	defer programMu.Unlock()
	if program == nil {
		panic("program: Program is nil")
	}
	if _, dup := programs[name]; dup {
		panic("program: Register called twice for program " + name)
	}
	programs[name] = program
}

func Resolve(name string) Program {
	programMu.RLock()
	defer programMu.RUnlock()

	program, ok := programs[name]
	if !ok {
		panic("Program has not been registered: " + name)
	}

	return program
}
