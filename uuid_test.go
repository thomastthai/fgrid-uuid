package uuid

import (
	"encoding"
	"fmt"
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

func ExampleUUID_String() {
	fmt.Printf("NIL-UUID: %s", NIL.String())
	// Output:
	// NIL-UUID: 00000000-0000-0000-0000-000000000000
}

func TestParseUUID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid hyphenated", "00000000-0000-0000-0000-000000000000", false},
		{"valid no hyphens", "00000000000000000000000000000000", false},
		{"valid mixed case", "ABCDEF12-3456-7890-ABCD-EF1234567890", false},
		{"invalid length short", "00000000-0000-0000-0000-00000000000", true},
		{"invalid length long", "00000000-0000-0000-0000-0000000000000", true},
		{"invalid characters", "gggggggg-0000-0000-0000-000000000000", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid, err := ParseUUID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && uuid == nil {
				t.Errorf("ParseUUID() returned nil UUID for valid input")
			}
		})
	}
}

func TestUUIDMarshalText(t *testing.T) {
	uuid := NIL
	text, err := uuid.MarshalText()
	if err != nil {
		t.Errorf("MarshalText() error = %v", err)
	}
	expected := "00000000-0000-0000-0000-000000000000"
	if string(text) != expected {
		t.Errorf("MarshalText() = %s, want %s", string(text), expected)
	}
}

func TestUUIDUnmarshalText(t *testing.T) {
	var uuid UUID
	text := []byte("00000000-0000-0000-0000-000000000000")
	err := uuid.UnmarshalText(text)
	if err != nil {
		t.Errorf("UnmarshalText() error = %v", err)
	}
	if uuid != *NIL {
		t.Errorf("UnmarshalText() = %v, want %v", uuid, *NIL)
	}
}

func TestUUIDImplementsInterfaces(t *testing.T) {
	var uuid UUID
	
	// Test that UUID implements the expected interfaces
	var _ fmt.Stringer = &uuid
	var _ encoding.TextMarshaler = &uuid
	var _ encoding.TextUnmarshaler = &uuid
}

func TestParseUUIDRoundTrip(t *testing.T) {
	// Test with V1
	v1 := NewV1()
	parsed, err := ParseUUID(v1.String())
	if err != nil {
		t.Errorf("ParseUUID() error = %v", err)
	}
	if *parsed != *v1 {
		t.Errorf("ParseUUID round trip failed: got %v, want %v", parsed, v1)
	}
	
	// Test with V4
	v4 := NewV4()
	parsed, err = ParseUUID(v4.String())
	if err != nil {
		t.Errorf("ParseUUID() error = %v", err)
	}
	if *parsed != *v4 {
		t.Errorf("ParseUUID round trip failed: got %v, want %v", parsed, v4)
	}
}
