// Package uuid provides functionality for generating and working with UUIDs
// according to RFC 4122 and RFC 9562.
//
// This package supports UUID versions 1, 3, 4, 5, and 7:
//   - Version 1: timestamp and MAC address based
//   - Version 3: namespace and name with MD5 hash
//   - Version 4: random or pseudo-random numbers
//   - Version 5: namespace and name with SHA-1 hash
//   - Version 7: timestamp with random data (RFC 9562)
package uuid

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// The UUID represents Universally Unique IDentifier (which is 128 bit long).
// A UUID is a 16-byte array that uniquely identifies an object or entity.
type UUID [16]byte

var (
	// NIL is defined in RFC 4122 section 4.1.7.
	// The nil UUID is special form of UUID that is specified to have all 128 bits set to zero.
	NIL = &UUID{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	// NameSpaceDNS assume name to be a fully-qualified domain name.
	// Declared in RFC 4122 Appendix C.
	NameSpaceDNS = &UUID{
		0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1,
		0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8,
	}
	// NameSpaceURL assume name to be a URL.
	// Declared in RFC 4122 Appendix C.
	NameSpaceURL = &UUID{
		0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1,
		0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8,
	}
	// NameSpaceOID assume name to be an ISO OID.
	// Declared in RFC 4122 Appendix C.
	NameSpaceOID = &UUID{
		0x6b, 0xa7, 0xb8, 0x12, 0x9d, 0xad, 0x11, 0xd1,
		0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8,
	}
	// NameSpaceX500 assume name to be a X.500 DN (in DER or a text output format).
	// Declared in RFC 4122 Appendix C.
	NameSpaceX500 = &UUID{
		0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1,
		0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8,
	}
)

// Version returns the version number of the UUID.
// The version indicates which algorithm was used to generate the UUID.
func (u *UUID) Version() int {
	return int(binary.BigEndian.Uint16(u[6:8]) >> 12)
}

// String returns the human readable form of the UUID.
// The format is xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx where x is a hexadecimal digit.
func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// ParseUUID parses a UUID string in the standard format and returns a UUID.
// The string must be in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// where x is a hexadecimal digit (case insensitive).
func ParseUUID(s string) (*UUID, error) {
	// Remove any whitespace and convert to lowercase
	s = strings.ToLower(strings.TrimSpace(s))
	
	// Check the length and format
	if len(s) != 36 {
		return nil, errors.New("invalid UUID format: incorrect length")
	}
	
	// Check for hyphens in the correct positions
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return nil, errors.New("invalid UUID format: missing hyphens")
	}
	
	// Remove hyphens for hex parsing
	hex := strings.Replace(s, "-", "", -1)
	if len(hex) != 32 {
		return nil, errors.New("invalid UUID format: incorrect hex length")
	}
	
	// Parse each hex digit
	var uuid UUID
	for i := 0; i < 32; i += 2 {
		var b byte
		for j := 0; j < 2; j++ {
			c := hex[i+j]
			var n byte
			if c >= '0' && c <= '9' {
				n = c - '0'
			} else if c >= 'a' && c <= 'f' {
				n = c - 'a' + 10
			} else {
				return nil, fmt.Errorf("invalid UUID format: invalid hex character '%c'", c)
			}
			if j == 0 {
				b = n << 4
			} else {
				b |= n
			}
		}
		uuid[i/2] = b
	}
	
	return &uuid, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It allows UUIDs to be unmarshaled from JSON and other text formats.
func (u *UUID) UnmarshalText(text []byte) error {
	parsed, err := ParseUUID(string(text))
	if err != nil {
		return err
	}
	*u = *parsed
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// It allows UUIDs to be marshaled to JSON and other text formats.
func (u *UUID) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// MarshalJSON implements the json.Marshaler interface.
// It ensures UUIDs are marshaled as strings in JSON format.
func (u *UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It allows UUIDs to be unmarshaled from JSON strings.
func (u *UUID) UnmarshalJSON(data []byte) error {
	// Remove quotes from JSON string
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("invalid JSON format for UUID")
	}
	
	// Parse the UUID string
	parsed, err := ParseUUID(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*u = *parsed
	return nil
}

// variantRFC4122 sets the variant bits according to RFC 4122.
// This is an internal method used during UUID generation.
func (u *UUID) variantRFC4122() {
	u[8] = (u[8] & 0x3f) | 0x80
}
