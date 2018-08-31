package epub

import (
	"crypto/rand"
	"fmt"
	"io"
)

// NewUUID returns the canonical string representation of a UUID:
//  xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func NewUUID() string {
	var uuid [16]byte
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // set version byte
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // set high order byte 0b10{8,9,a,b}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
