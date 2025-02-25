package cmdutil

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	inputFile   *os.File = os.Stdin
	inputBuffer *bufio.Reader
)

// ReadLine - does stuff
func ReadLine(inputReader *bufio.Reader) string {
	line, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSpace(string(line))
}

// Silence - does stuff
func Silence() {
	runCommand(exec.Command("stty", "-echo"))
}

// Unsilence - does stuff
func Unsilence() {
	runCommand(exec.Command("stty", "echo"))
}

func runCommand(command *exec.Cmd) {
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Run()
}

func buffer() *bufio.Reader {
	if inputBuffer == nil {
		inputBuffer = bufio.NewReader(inputFile)
	}
	return inputBuffer
}
