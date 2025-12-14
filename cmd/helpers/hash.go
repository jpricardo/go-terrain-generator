package helpers

import (
	"encoding/binary"
	"hash/maphash"
)

func HashInt(n int) uint8 {
	var h maphash.Hash
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(n))

	h.Write(b)
	return uint8(h.Sum64())
}
