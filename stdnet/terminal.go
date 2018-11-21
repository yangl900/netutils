package main

import (
	"log"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	dll            = syscall.MustLoadDLL("kernel32")
	setConsoleMode = dll.MustFindProc("SetConsoleMode")
	memory         uint32
)

func setInputConsoleMode(h syscall.Handle, m uint32) error {
	r, _, err := setConsoleMode.Call(uintptr(h), uintptr(m))
	if r == 0 {
		return err
	}
	return nil
}

type terminalState struct {
	state *terminal.State
}

func isTerminal(fd uintptr) bool {
	return terminal.IsTerminal(int(fd))
}

func setConsoleRaw() {
	h := syscall.Handle(os.Stdin.Fd())

	if err := syscall.GetConsoleMode(h, &memory); err != nil {
		log.Fatal(err)
	}
	if err := setInputConsoleMode(h, 0); err != nil {
		log.Fatal(err)
	}
}

func restoreConsole() {
	h := syscall.Handle(os.Stdin.Fd())
	setInputConsoleMode(h, memory)
}

func makeRaw(fd uintptr) (*terminalState, error) {
	state, err := terminal.MakeRaw(int(fd))

	return &terminalState{
		state: state,
	}, err
}

func restore(fd uintptr, oldState *terminalState) error {
	return terminal.Restore(int(fd), oldState.state)
}
