package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type PathNotProvidedError struct{}

func (*PathNotProvidedError) Error() string {
	return "path not provided"
}

type Shell struct {
	reader     *bufio.Reader
	input      string
	enablePath bool
}

func (s *Shell) getPath() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return path
}

func (s *Shell) init(enablePath bool) {
	s.reader = bufio.NewReader(os.Stdin)
	s.enablePath = enablePath
}

func (s *Shell) getInput() {
	if s.enablePath {
		fmt.Print(s.getPath())
	}
	fmt.Print("> ")
	input, err := s.reader.ReadString('\n')
	if err != nil {
		os.Exit(0)
	}
	s.input = input
}

func (s *Shell) run() {
	for {
		s.getInput()
		if err := s.execInput(); err != nil {
			switch err.(type) {
			case *exec.Error, *exec.ExitError, *os.PathError, *PathNotProvidedError:
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func (s *Shell) execInput() error {
	input := strings.TrimSpace(s.input)
	args := strings.Split(input, " ")
	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return &PathNotProvidedError{}
		}
		return os.Chdir(args[1])
	case "shell":
		fmt.Println("kshell is a simple shell written in go by Kushal Chandar.")
		return nil
	case "path":
		s.enablePath = !s.enablePath
		return nil
	case "exit":
		os.Exit(0)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func main() {
	var shell Shell
	shell.init(true)
	shell.run()
}
