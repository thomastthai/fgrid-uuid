package uuid

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
)

// The UUID represents Universally Unique IDentifier (which is 128 bit long).
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

// Version of the UUID represents a kind of subtype specifier.
func (u *UUID) Version() int {
	return int(binary.BigEndian.Uint16(u[6:8]) >> 12)
}

// String returns the human readable form of the UUID.
func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// MarshalText implements encoding.TextMarshaler interface.
// Returns the string representation of the UUID.
func (u *UUID) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface.
// Parses a UUID from its string representation.
func (u *UUID) UnmarshalText(text []byte) error {
	parsed, err := ParseUUID(string(text))
	if err != nil {
		return err
	}
	*u = *parsed
	return nil
}

// ParseUUID parses a UUID from its string representation.
// Accepts both hyphenated (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx) and 
// non-hyphenated (xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx) formats.
// Returns an error if the input is not a valid UUID string.
func ParseUUID(s string) (*UUID, error) {
	// Remove hyphens and convert to lowercase
	s = strings.ReplaceAll(strings.ToLower(s), "-", "")
	
	// Validate length
	if len(s) != 32 {
		return nil, fmt.Errorf("invalid UUID length: expected 32 hex characters, got %d", len(s))
	}
	
	// Validate hex characters
	hexPattern := regexp.MustCompile("^[0-9a-f]{32}$")
	if !hexPattern.MatchString(s) {
		return nil, fmt.Errorf("invalid UUID format: contains non-hex characters")
	}
	
	var uuid UUID
	for i := 0; i < 16; i++ {
		// Parse two hex characters at a time
		var b byte
		if _, err := fmt.Sscanf(s[i*2:i*2+2], "%02x", &b); err != nil {
			return nil, fmt.Errorf("invalid UUID format: %v", err)
		}
		uuid[i] = b
	}
	
	return &uuid, nil
}

func (u *UUID) variantRFC4122() {
	u[8] = (u[8] & 0x3f) | 0x80
}
