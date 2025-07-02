// Package uuid implements UUID version 3 generation.
// This file contains the implementation for namespace and name based UUIDs using MD5 hash.
package uuid

import (
	"crypto/md5"
	"hash"
)

// NewV3 creates a new UUID with version 3 as described in RFC 4122.
// Version 3 UUIDs are based on a namespace UUID and a name, using MD5 hash calculation.
// The resulting UUID is deterministic - the same namespace and name will always
// produce the same UUID.
//
// Parameters:
//   - namespace: A UUID that represents the namespace
//   - name: A byte slice containing the name to hash
//
// Returns a pointer to the generated UUID.
func NewV3(namespace *UUID, name []byte) *UUID {
	uuid := newByHash(md5.New(), namespace, name)
	uuid[6] = (uuid[6] & 0x0f) | 0x30  // Set version to 3
	return uuid
}

// newByHash is a helper function that generates a UUID from a hash function,
// namespace UUID, and name. This function is shared between UUID v3 and v5
// implementations, differing only in the hash algorithm used.
func newByHash(hash hash.Hash, namespace *UUID, name []byte) *UUID {
	// Hash the namespace UUID bytes followed by the name bytes
	hash.Write(namespace[:])
	hash.Write(name[:])

	var uuid UUID
	// Use first 16 bytes of the hash as the UUID
	copy(uuid[:], hash.Sum(nil)[:16])
	
	// Set variant bits according to RFC 4122
	uuid.variantRFC4122()
	return &uuid
}
