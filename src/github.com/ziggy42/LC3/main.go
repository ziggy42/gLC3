package main

import (
	"fmt"
	"os"
)

// Registers
const (
	R_R0    = iota
	R_R1    = iota
	R_R2    = iota
	R_R3    = iota
	R_R4    = iota
	R_R5    = iota
	R_R6    = iota
	R_R7    = iota
	R_PC    = iota
	R_COND  = iota
	R_COUNT = iota
)

// OPCodes
const (
	OP_BR   = iota // branch
	OP_ADD  = iota // add
	OP_LD   = iota // load
	OP_ST   = iota // store
	OP_JSR  = iota // jump register
	OP_AND  = iota // bitwise and
	OP_LDR  = iota // load register
	OP_STR  = iota // store register
	OP_RTI  = iota // unused
	OP_NOT  = iota // bitwise not
	OP_LDI  = iota // load indirect
	OP_STI  = iota // store indirect
	OP_JMP  = iota // jump
	OP_RES  = iota // reserved (unused)
	OP_LEA  = iota // load effective address
	OP_TRAP = iota // execute trap
)

const (
	FL_POS = 1 << 0
	FL_ZRO = 1 << 1
	FL_NEG = 1 << 2
)

const (
	TRAP_GETC  = 0x20 // get character from keyboard
	TRAP_OUT   = 0x21 // output a character
	TRAP_PUTS  = 0x22 // output a word string
	TRAP_IN    = 0x23 // input a string
	TRAP_PUTSP = 0x24 // output a byte string
	TRAP_HALT  = 0x25 // halt the program
)

const (
	MR_KBSR = 0xFE00 // keyboard status
	MR_KBDR = 0xFE02 // keyboard data
)

const (
	MODE_REGISTER  = iota
	MODE_IMMEDIATE = iota
)

// PC_START sets the PC starting position. 0x3000 is the default position
const PC_START = 0x3000

var memory [^uint16(0)]uint16
var registers [R_COUNT]uint16

func memWrite(address uint16, value uint16) {
	memory[address] = value
}

func memRead(address uint16) uint16 {
	if address == MR_KBSR {
		if IsKeyPressed() {
			memory[MR_KBSR] = (1 << 15)
			c, err := GetChar()
			if err != nil {
				panic(err)
			}

			memory[MR_KBDR] = c
		} else {
			memory[MR_KBSR] = 0
		}
	}

	return memory[address]
}

func updateFlags(r uint16) {
	if registers[r] == 0 {
		registers[R_COND] = FL_ZRO
	} else if registers[r]<<15 == 1 { // Right-most bit is 1 for negative numbers
		registers[R_COND] = FL_NEG
	} else {
		registers[R_COND] = FL_POS
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Missing argument!")
		return
	}

	if err := Load(&memory, args[0]); err != nil {
		panic(err)
	}

	registers[R_PC] = PC_START
	running := true
	for running {
		instruction := memory[registers[R_PC]]
		registers[R_PC]++

		switch op := instruction >> 12; op {
		case OP_BR:
			pcOffset := SignExtend(instruction&0x1ff, 9)
			flag := (instruction >> 9) & 0x7

			if (flag & registers[R_COND]) != 0 {
				registers[R_PC] += pcOffset
			}
		case OP_ADD:
			dr := (instruction >> 9) & 0x7
			sr1 := (instruction >> 6) & 0x7
			mode := (instruction >> 5) & 0x1

			if mode == MODE_IMMEDIATE {
				immediate := SignExtend(instruction&0x1F, 5)
				registers[dr] = registers[sr1] + immediate
			} else {
				sr2 := instruction & 0x7
				registers[dr] = registers[sr1] + registers[sr2]
			}

			updateFlags(dr)
		case OP_LD:
			dr := (instruction >> 9) & 0x7
			pcOffset := SignExtend(instruction&0x1ff, 9)

			registers[dr] = memRead(registers[R_PC] + pcOffset)
			updateFlags(dr)
		case OP_ST:
			sr := (instruction >> 9) & 0x7
			pcOffset := SignExtend(instruction&0x1ff, 9)
			memWrite(registers[R_PC]+pcOffset, registers[sr])
		case OP_JSR:
			r := (instruction >> 6) & 0x7
			longPcOffset := SignExtend(instruction&0x7ff, 11)
			longFlag := (instruction >> 11) & 1

			registers[R_R7] = registers[R_PC]
			if longFlag != 0 {
				registers[R_PC] += longPcOffset
			} else {
				registers[R_PC] = registers[r]
			}
		case OP_AND:
			dr := (instruction >> 9) & 0x7
			sr1 := (instruction >> 6) & 0x7
			mode := (instruction >> 5) & 0x1

			if mode == MODE_IMMEDIATE {
				imm5 := SignExtend(instruction&0x1F, 5)
				registers[dr] = registers[sr1] & imm5
			} else {
				sr2 := instruction & 0x7
				registers[dr] = registers[sr1] & registers[sr2]
			}

			updateFlags(dr)
		case OP_LDR:
			dr := (instruction >> 9) & 0x7
			sr1 := (instruction >> 6) & 0x7
			offset := SignExtend(instruction&0x3F, 6)
			registers[dr] = memRead(registers[sr1] + offset)

			updateFlags(dr)
		case OP_STR:
			sr1 := (instruction >> 9) & 0x7
			sr2 := (instruction >> 6) & 0x7
			offset := SignExtend(instruction&0x3F, 6)
			memWrite(registers[sr2]+offset, registers[sr1])
		case OP_RTI:
			panic("Unknown opcode")
		case OP_NOT:
			dr := (instruction >> 9) & 0x7
			sr1 := (instruction >> 6) & 0x7

			registers[dr] = ^registers[sr1]
			updateFlags(dr)
		case OP_LDI:
			dr := (instruction >> 9) & 0x7
			pcOffset := SignExtend(instruction&0x1ff, 9)

			registers[dr] = memRead(memRead(registers[R_PC] + pcOffset))
			updateFlags(dr)
		case OP_STI:
			sr1 := (instruction >> 9) & 0x7
			pcOffset := SignExtend(instruction&0x1ff, 9)
			memWrite(memRead(registers[R_PC]+pcOffset), registers[sr1])
		case OP_JMP:
			sr1 := (instruction >> 6) & 0x7
			registers[R_PC] = registers[sr1]
		case OP_RES:
			panic("Unknown opcode")
		case OP_LEA:
			dr := (instruction >> 9) & 0x7
			pcOffset := SignExtend(instruction&0x1ff, 9)

			registers[dr] = registers[R_PC] + pcOffset
			updateFlags(dr)
		case OP_TRAP:
			switch instruction & 0xFF {
			case TRAP_GETC:
				c, err := GetChar()
				if err != nil {
					panic(err)
				}
				registers[R_R0] = c
			case TRAP_OUT:
				fmt.Printf("%c", rune(registers[R_R0]))
			case TRAP_PUTS:
				i := registers[R_R0]
				for ; memory[i] != 0; i++ {
					fmt.Printf("%c", rune(memory[i]))
				}
			case TRAP_IN:
				fmt.Print("Enter a character: ")
				c, err := GetChar()
				if err != nil {
					panic(err)
				}
				registers[R_R0] = c
			case TRAP_PUTSP:
				for i := registers[R_R0]; memory[i] != 0; i++ {
					r1 := rune(memory[i] & 0xFF)
					fmt.Printf("%c", r1)
					r2 := rune(memory[i] >> 8)
					if r2 != 0 {
						fmt.Printf("%c", r2)
					}
				}
			case TRAP_HALT:
				running = false
			}
		default:
			panic("Unknown opcode")
		}
	}
}
