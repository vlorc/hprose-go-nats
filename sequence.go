package hprose_go_nats

import "strings"

type Sequence uint64

var hexTable [16]byte = [16]byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46}
var byteTable map[byte]byte = map[byte]byte{0x30: 0, 0x31: 1, 0x32: 2, 0x33: 3, 0x34: 4, 0x35: 5, 0x36: 6, 0x37: 7, 0x38: 8, 0x39: 9, 0x41: 10, 0x42: 11, 0x43: 12, 0x44: 13, 0x45: 14, 0x46: 15}

func NewSequence(s string) Sequence {
	pos := strings.LastIndexByte(s, byte('.'))
	seq := uint64(0)
	for _, v := range s[pos+1:] {
		seq <<= 4
		seq += uint64(byteTable[byte(v)])
	}
	return Sequence(seq)
}

func (s Sequence) String() string {
	if 0 == s {
		return ".0"
	}
	var buf [20]byte
	var pos = 19
	for ; s > 0; s >>= 4 {
		buf[pos] = hexTable[s&15]
		pos--
	}
	buf[pos] = byte('.')
	return string(buf[pos:])
}
