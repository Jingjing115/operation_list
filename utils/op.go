package utils

import (
	"encoding/hex"
	"encoding/binary"
	"log"
)

func ToAddr(addr string) uint32 {
	data, err := hex.DecodeString(addr)
	if err != nil {
		log.Println(err)
		return 0
	}
	return binary.LittleEndian.Uint32(data)
}

