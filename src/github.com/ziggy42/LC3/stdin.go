package main

import (
	"os"
	"os/exec"
)

// GetChar reads one character from Stdin as an uint16
func GetChar() (uint16, error) {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	b := make([]byte, 1)
	os.Stdin.Read(b)
	return uint16(b[0]), nil
}

// IsKeyPressed checks if a key was pressed
func IsKeyPressed() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Size() > 0
}
