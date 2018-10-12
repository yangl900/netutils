package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	cmd := exec.Command("/bin/bash", "-c", "stty -echo; /ssh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println("Error when starting process: ", err.Error())
	}

	cmd.Wait()
}
