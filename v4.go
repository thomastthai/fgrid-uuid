// Package uuid implements UUID version 4 generation.
// This file contains the implementation for random-based UUIDs.
package uuid

import "crypto/rand"

// NewV4 creates a new UUID with version 4 as described in RFC 4122.
// Version 4 UUIDs are based on random or pseudo-random numbers.
// This is the most commonly used UUID version due to its simplicity
// and strong uniqueness properties.
//
// The UUID is generated using cryptographically secure random bytes,
// ensuring high entropy and unpredictability.
//
// Returns a pointer to the generated UUID.
func NewV4() *UUID {
	// Generate 16 random bytes
	buf := make([]byte, 16)
	rand.Read(buf)
	
	// Set version to 4 (0100 in binary) in the upper 4 bits of byte 6
	buf[6] = (buf[6] & 0x0f) | 0x40
	
	var uuid UUID
	copy(uuid[:], buf[:])
	
	// Set variant bits according to RFC 4122
	uuid.variantRFC4122()
	return &uuid
}
