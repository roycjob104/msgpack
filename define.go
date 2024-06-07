package msgpack

const (
	Max8Bit  = (1 << 8) - 1
	Max16Bit = (1 << 16) - 1
	Max32Bit = (1 << 32) - 1
)

var (
	PosFixedNumHigh byte = 0x7f
	NegFixedNumLow  byte = 0xe0

	Nil byte = 0xc0

	False byte = 0xc2
	True  byte = 0xc3

	Float64 byte = 0xcb

	FixedStrLow  byte = 0xa0
	FixedStrHigh byte = 0xbf
	Str8         byte = 0xd9
	Str16        byte = 0xda
	Str32        byte = 0xdb

	FixedArrayLow  byte = 0x90
	FixedArrayHigh byte = 0x9f
	Array16        byte = 0xdc
	Array32        byte = 0xdd

	FixedMapLow  byte = 0x80
	FixedMapHigh byte = 0x8f
	Map16        byte = 0xde
	Map32        byte = 0xdf
)
