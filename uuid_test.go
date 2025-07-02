package uuid

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	uuid := UUID{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00,
		0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	uuid[6] = 0x10
	if uuid.Version() != 1 {
		t.Errorf("invalid version %d - expected 1", uuid.Version())
	}
	uuid[6] = 0x20
	if uuid.Version() != 2 {
		t.Errorf("invalid version %d - expected 2", uuid.Version())
	}
	uuid[6] = 0x30
	if uuid.Version() != 3 {
		t.Errorf("invalid version %d - expected 3", uuid.Version())
	}
	uuid[6] = 0x40
	if uuid.Version() != 4 {
		t.Errorf("invalid version %d - expected 4", uuid.Version())
	}
	uuid[6] = 0x50
	if uuid.Version() != 5 {
		t.Errorf("invalid version %d - expected 5", uuid.Version())
	}
	uuid[6] = 0x70
	if uuid.Version() != 7 {
		t.Errorf("invalid version %d - expected 7", uuid.Version())
	}
}

// TestParseUUID tests the ParseUUID function with various valid and invalid inputs
func TestParseUUID(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  *UUID
	}{
		{
			name:      "nil UUID",
			input:     "00000000-0000-0000-0000-000000000000",
			shouldErr: false,
			expected:  NIL,
		},
		{
			name:      "valid UUID with uppercase",
			input:     "F47AC10B-58CC-4372-A567-0E02B2C3D479",
			shouldErr: false,
			expected: &UUID{
				0xf4, 0x7a, 0xc1, 0x0b, 0x58, 0xcc, 0x43, 0x72,
				0xa5, 0x67, 0x0e, 0x02, 0xb2, 0xc3, 0xd4, 0x79,
			},
		},
		{
			name:      "valid UUID with lowercase",
			input:     "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			shouldErr: false,
			expected: &UUID{
				0xf4, 0x7a, 0xc1, 0x0b, 0x58, 0xcc, 0x43, 0x72,
				0xa5, 0x67, 0x0e, 0x02, 0xb2, 0xc3, 0xd4, 0x79,
			},
		},
		{
			name:      "valid UUID with mixed case",
			input:     "F47ac10B-58CC-4372-a567-0E02b2c3D479",
			shouldErr: false,
			expected: &UUID{
				0xf4, 0x7a, 0xc1, 0x0b, 0x58, 0xcc, 0x43, 0x72,
				0xa5, 0x67, 0x0e, 0x02, 0xb2, 0xc3, 0xd4, 0x79,
			},
		},
		{
			name:      "valid UUID with whitespace",
			input:     "  f47ac10b-58cc-4372-a567-0e02b2c3d479  ",
			shouldErr: false,
			expected: &UUID{
				0xf4, 0x7a, 0xc1, 0x0b, 0x58, 0xcc, 0x43, 0x72,
				0xa5, 0x67, 0x0e, 0x02, 0xb2, 0xc3, 0xd4, 0x79,
			},
		},
		{
			name:      "too short",
			input:     "f47ac10b-58cc-4372-a567-0e02b2c3d47",
			shouldErr: true,
		},
		{
			name:      "too long",
			input:     "f47ac10b-58cc-4372-a567-0e02b2c3d4790",
			shouldErr: true,
		},
		{
			name:      "missing hyphens",
			input:     "f47ac10b58cc4372a5670e02b2c3d479",
			shouldErr: true,
		},
		{
			name:      "wrong hyphen positions",
			input:     "f47ac10b5-8cc-4372-a567-0e02b2c3d479",
			shouldErr: true,
		},
		{
			name:      "invalid hex character",
			input:     "f47ac10g-58cc-4372-a567-0e02b2c3d479",
			shouldErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			shouldErr: true,
		},
		{
			name:      "only hyphens",
			input:     "--------",
			shouldErr: true,
		},
		{
			name:      "only whitespace",
			input:     "   ",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseUUID(tt.input)
			
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error for input %q, but got none", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
				return
			}
			
			if result == nil {
				t.Errorf("expected UUID result for input %q, but got nil", tt.input)
				return
			}
			
			if *result != *tt.expected {
				t.Errorf("for input %q:\nexpected: %v\ngot:      %v", tt.input, tt.expected, result)
			}
		})
	}
}

// TestUnmarshalText tests the UnmarshalText method
func TestUnmarshalText(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		expected  UUID
	}{
		{
			name:      "valid UUID",
			input:     "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			shouldErr: false,
			expected: UUID{
				0xf4, 0x7a, 0xc1, 0x0b, 0x58, 0xcc, 0x43, 0x72,
				0xa5, 0x67, 0x0e, 0x02, 0xb2, 0xc3, 0xd4, 0x79,
			},
		},
		{
			name:      "nil UUID",
			input:     "00000000-0000-0000-0000-000000000000",
			shouldErr: false,
			expected:  *NIL,
		},
		{
			name:      "invalid UUID",
			input:     "invalid-uuid",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uuid UUID
			err := uuid.UnmarshalText([]byte(tt.input))
			
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error for input %q, but got none", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
				return
			}
			
			if uuid != tt.expected {
				t.Errorf("for input %q:\nexpected: %v\ngot:      %v", tt.input, tt.expected, uuid)
			}
		})
	}
}

// TestMarshalText tests the MarshalText method
func TestMarshalText(t *testing.T) {
	uuid := UUID{
		0xf4, 0x7a, 0xc1, 0x0b, 0x58, 0xcc, 0x43, 0x72,
		0xa5, 0x67, 0x0e, 0x02, 0xb2, 0xc3, 0xd4, 0x79,
	}
	
	result, err := uuid.MarshalText()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	expected := "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	if string(result) != expected {
		t.Errorf("expected %q, got %q", expected, string(result))
	}
}

// TestJSONMarshaling tests JSON marshaling and unmarshaling
func TestJSONMarshaling(t *testing.T) {
	original := UUID{
		0xf4, 0x7a, 0xc1, 0x0b, 0x58, 0xcc, 0x43, 0x72,
		0xa5, 0x67, 0x0e, 0x02, 0xb2, 0xc3, 0xd4, 0x79,
	}
	
	// Marshal to JSON using pointer (required for array types)
	jsonBytes, err := json.Marshal(&original)
	if err != nil {
		t.Errorf("unexpected error marshaling to JSON: %v", err)
	}
	
	expected := `"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	if string(jsonBytes) != expected {
		t.Errorf("expected JSON %s, got %s", expected, string(jsonBytes))
	}
	
	// Unmarshal from JSON
	var unmarshaled UUID
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	if err != nil {
		t.Errorf("unexpected error unmarshaling from JSON: %v", err)
	}
	
	if unmarshaled != original {
		t.Errorf("expected %v, got %v", original, unmarshaled)
	}
}

// TestRoundTripParsing tests that String() and ParseUUID() are inverse operations
func TestRoundTripParsing(t *testing.T) {
	// Test with various UUIDs
	testUUIDs := []*UUID{
		NIL,
		NewV1(),
		NewV3(NameSpaceDNS, []byte("example.com")),
		NewV4(),
		NewV5(NameSpaceDNS, []byte("example.com")),
		NewV7(),
	}
	
	for i, original := range testUUIDs {
		t.Run(fmt.Sprintf("UUID_%d", i), func(t *testing.T) {
			// Convert to string and back
			str := original.String()
			parsed, err := ParseUUID(str)
			if err != nil {
				t.Errorf("unexpected error parsing %q: %v", str, err)
			}
			
			if *parsed != *original {
				t.Errorf("round trip failed:\noriginal: %v\nparsed:   %v", original, parsed)
			}
		})
	}
}

// TestUUIDCreationEdgeCases tests edge cases in UUID creation
func TestUUIDCreationEdgeCases(t *testing.T) {
	t.Run("V1_Multiple_Calls", func(t *testing.T) {
		// Generate multiple V1 UUIDs rapidly to test uniqueness and timestamp handling
		uuids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			uuid := NewV1()
			str := uuid.String()
			if uuids[str] {
				t.Errorf("duplicate V1 UUID generated: %s", str)
			}
			uuids[str] = true
			
			if uuid.Version() != 1 {
				t.Errorf("expected version 1, got %d", uuid.Version())
			}
		}
	})
	
	t.Run("V3_Same_Input", func(t *testing.T) {
		// V3 UUIDs should be deterministic
		uuid1 := NewV3(NameSpaceDNS, []byte("test"))
		uuid2 := NewV3(NameSpaceDNS, []byte("test"))
		
		if *uuid1 != *uuid2 {
			t.Errorf("V3 UUIDs with same input should be identical:\nUUID1: %s\nUUID2: %s", uuid1.String(), uuid2.String())
		}
		
		if uuid1.Version() != 3 {
			t.Errorf("expected version 3, got %d", uuid1.Version())
		}
	})
	
	t.Run("V4_Uniqueness", func(t *testing.T) {
		// V4 UUIDs should be unique
		uuids := make(map[string]bool)
		for i := 0; i < 1000; i++ {
			uuid := NewV4()
			str := uuid.String()
			if uuids[str] {
				t.Errorf("duplicate V4 UUID generated: %s", str)
			}
			uuids[str] = true
			
			if uuid.Version() != 4 {
				t.Errorf("expected version 4, got %d", uuid.Version())
			}
		}
	})
	
	t.Run("V5_Same_Input", func(t *testing.T) {
		// V5 UUIDs should be deterministic
		uuid1 := NewV5(NameSpaceDNS, []byte("test"))
		uuid2 := NewV5(NameSpaceDNS, []byte("test"))
		
		if *uuid1 != *uuid2 {
			t.Errorf("V5 UUIDs with same input should be identical:\nUUID1: %s\nUUID2: %s", uuid1.String(), uuid2.String())
		}
		
		if uuid1.Version() != 5 {
			t.Errorf("expected version 5, got %d", uuid1.Version())
		}
	})
	
	t.Run("V7_Ordering", func(t *testing.T) {
		// V7 UUIDs should generally be ordered by creation time
		var previousUUID *UUID
		for i := 0; i < 10; i++ {
			uuid := NewV7()
			if uuid.Version() != 7 {
				t.Errorf("expected version 7, got %d", uuid.Version())
			}
			
			if previousUUID != nil {
				// Compare timestamp portions (first 6 bytes)
				// Due to millisecond precision, some UUIDs might have the same timestamp
				previousTimestamp := string(previousUUID[:6])
				currentTimestamp := string(uuid[:6])
				if strings.Compare(previousTimestamp, currentTimestamp) > 0 {
					t.Errorf("V7 UUID ordering violation: later UUID has earlier timestamp")
				}
			}
			previousUUID = uuid
		}
	})
	
	t.Run("Namespace_UUIDs", func(t *testing.T) {
		// Test standard namespace UUIDs
		namespaces := []*UUID{NameSpaceDNS, NameSpaceURL, NameSpaceOID, NameSpaceX500}
		for i, ns := range namespaces {
			if ns == nil {
				t.Errorf("namespace %d is nil", i)
			}
			
			// All standard namespaces should be valid UUIDs
			str := ns.String()
			if len(str) != 36 {
				t.Errorf("namespace %d has invalid string representation: %s", i, str)
			}
		}
	})
}

func ExampleUUID_String() {
	fmt.Printf("NIL-UUID: %s", NIL.String())
	// Output:
	// NIL-UUID: 00000000-0000-0000-0000-000000000000
}

func ExampleParseUUID() {
	uuid, err := ParseUUID("f47ac10b-58cc-4372-a567-0e02b2c3d479")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Parsed UUID: %s\n", uuid.String())
	fmt.Printf("Version: %d\n", uuid.Version())
	// Output:
	// Parsed UUID: f47ac10b-58cc-4372-a567-0e02b2c3d479
	// Version: 4
}
