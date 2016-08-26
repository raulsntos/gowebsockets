package gowebsockets

import (
	"crypto/rand"
	"encoding/hex"
)

// UUID is a unique identifier.
// You won't use this type directly.
type UUID [16]byte

func generateUUID() string {
	u := UUID{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(u[:])
}
