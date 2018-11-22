package main

import (
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	state terminalState
)

type terminalState struct {
	state *terminal.State
}

func isTerminal(fd uintptr) bool {
	return terminal.IsTerminal(int(fd))
}

func setConsoleRaw() {
	s, err := terminal.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		log.Fatal(err)
	}

	state.state = s
}

func restoreConsole() {
	terminal.Restore(int(os.Stdin.Fd()), state.state)
}
