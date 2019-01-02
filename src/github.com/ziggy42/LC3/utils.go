package main

// SignExtend extends the sign of x of bitCount bits
func SignExtend(x uint16, bitCount uint) uint16 {
	if ((x >> (bitCount - 1)) & 1) > 0 {
		x |= (0xFFFF << bitCount)
	}
	return x
}
