// Copyright 2016 Raul Santos Lebrato
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
