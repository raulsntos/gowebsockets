package websockets

import "crypto/rand"

// UUID is a unique identifier
type UUID [16]byte

func generateUUID() string {
	u := UUID{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return string(u[:])
}
