package main

import (
	"bufio"
	"os"
)

// GetChar reads one character from Stdin as an uint16
func GetChar() (uint16, error) {
	r := bufio.NewReader(os.Stdin)
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint16(b), nil
}

// IsKeyPressed checks if a key was pressed
func IsKeyPressed() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Size() > 0
}
