package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

// NewV7 creates a new UUID with version 7 as described in RFC 9562.
// Version 7 is based on timestamp (milliseconds since Unix epoch) with random bits for uniqueness.
func NewV7() *UUID {
	var uuid UUID
	
	// Get current timestamp in milliseconds since Unix epoch
	timestampMs := uint64(time.Now().UnixMilli())
	
	// Fill first 6 bytes with timestamp (48 bits, big-endian)
	binary.BigEndian.PutUint64(uuid[0:8], timestampMs<<16)
	
	// Generate random bytes for the remaining portions
	randomBytes := make([]byte, 10)
	rand.Read(randomBytes)
	
	// Set the remaining random bits
	// Byte 6: 4 bits version + 4 bits random
	// Byte 7: 8 bits random  
	// Bytes 8-15: 64 bits random (byte 8 will be modified by variantRFC4122)
	copy(uuid[6:8], randomBytes[0:2])
	copy(uuid[8:16], randomBytes[2:10])
	
	// Set version to 7 (0111 in binary) in the upper 4 bits of byte 6
	uuid[6] = (uuid[6] & 0x0f) | 0x70
	
	// Set variant bits according to RFC 4122
	uuid.variantRFC4122()
	
	return &uuid
}