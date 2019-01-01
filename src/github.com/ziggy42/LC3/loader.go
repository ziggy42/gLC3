package main

import (
	"encoding/binary"
	"io/ioutil"
)

// Load loads a binary file located at the given path in the given buffer
func Load(memory *[^uint16(0)]uint16, path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	origin := binary.BigEndian.Uint16(b[:2])

	for i := 2; i < len(b); i += 2 {
		memory[origin] = binary.BigEndian.Uint16(b[i : i+2])
		origin++
	}

	return nil
}
