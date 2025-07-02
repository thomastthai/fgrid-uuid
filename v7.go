// Package uuid implements UUID version 7 generation.
// This file contains the implementation for timestamp-based UUIDs with random data as described in RFC 9562.
package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

// NewV7 creates a new UUID with version 7 as described in RFC 9562.
// Version 7 UUIDs are based on Unix timestamp in milliseconds with random bits for uniqueness.
// This version provides several advantages:
//   - Monotonic ordering by creation time
//   - High performance generation
//   - 48-bit timestamp (good until year 10,895 CE)
//   - Better database performance due to natural ordering
//
// The structure is:
//   - 48 bits: Unix timestamp in milliseconds (big-endian)
//   - 12 bits: Random data
//   - 4 bits: Version (0111 for version 7)
//   - 62 bits: Random data
//   - 2 bits: Variant (10 for RFC 4122)
//
// Returns a pointer to the generated UUID.
func NewV7() *UUID {
	var uuid UUID
	
	// Get current timestamp in milliseconds since Unix epoch
	timestampMs := uint64(time.Now().UnixMilli())
	
	// Fill first 6 bytes with timestamp (48 bits, big-endian)
	// We put the timestamp in the first 48 bits by shifting left 16 bits
	binary.BigEndian.PutUint64(uuid[0:8], timestampMs<<16)
	
	// Generate random bytes for the remaining portions
	randomBytes := make([]byte, 10)
	rand.Read(randomBytes)
	
	// Set the remaining random bits:
	// - Byte 6: 4 bits version + 4 bits random
	// - Byte 7: 8 bits random  
	// - Bytes 8-15: 64 bits random (byte 8 will be modified by variantRFC4122)
	copy(uuid[6:8], randomBytes[0:2])
	copy(uuid[8:16], randomBytes[2:10])
	
	// Set version to 7 (0111 in binary) in the upper 4 bits of byte 6
	uuid[6] = (uuid[6] & 0x0f) | 0x70
	
	// Set variant bits according to RFC 4122
	uuid.variantRFC4122()
	
	return &uuid
}