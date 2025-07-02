// Package uuid implements UUID version 5 generation.
// This file contains the implementation for namespace and name based UUIDs using SHA-1 hash.
package uuid

import "crypto/sha1"

// NewV5 creates a new UUID with version 5 as described in RFC 4122.
// Version 5 UUIDs are based on a namespace UUID and a name, using SHA-1 hash calculation.
// The resulting UUID is deterministic - the same namespace and name will always
// produce the same UUID. Version 5 is preferred over version 3 due to SHA-1's
// stronger cryptographic properties compared to MD5.
//
// Parameters:
//   - namespaceUUID: A UUID that represents the namespace
//   - name: A byte slice containing the name to hash
//
// Returns a pointer to the generated UUID.
func NewV5(namespaceUUID *UUID, name []byte) *UUID {
	uuid := newByHash(sha1.New(), namespaceUUID, name)
	uuid[6] = (uuid[6] & 0x0f) | 0x50  // Set version to 5
	return uuid
}

// NewNamespaceUUID creates a namespace UUID by using the namespace name in the NIL namespace.
// This is a different approach from the 4 "standard" namespace UUIDs which are time-based UUIDs (V1).
// This function provides a convenient way to create deterministic namespace UUIDs
// from string identifiers.
//
// Parameters:
//   - namespace: A string representing the namespace name
//
// Returns a pointer to the generated namespace UUID.
func NewNamespaceUUID(namespace string) *UUID {
	return NewV5(NIL, []byte(namespace))
}
