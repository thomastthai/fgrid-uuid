package uuid

import (
	"encoding/binary"
	"testing"
	"time"
)

func TestNewV7(t *testing.T) {
	uuid := NewV7()
	if uuid.Version() != 7 {
		t.Errorf("invalid version %d - expected 7", uuid.Version())
	}
	t.Logf("UUID V7: %s", uuid)
}

func TestV7Uniqueness(t *testing.T) {
	// Generate multiple UUIDs and ensure they are unique
	uuids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		uuid := NewV7()
		str := uuid.String()
		if uuids[str] {
			t.Errorf("duplicate UUID generated: %s", str)
		}
		uuids[str] = true
	}
}

func TestV7Ordering(t *testing.T) {
	// Generate UUIDs with small time delays and ensure they are ordered
	uuid1 := NewV7()
	time.Sleep(1 * time.Millisecond)
	uuid2 := NewV7()
	
	// Compare the timestamp portions (first 6 bytes)
	for i := 0; i < 6; i++ {
		if uuid1[i] < uuid2[i] {
			break // uuid1 is earlier, which is correct
		}
		if uuid1[i] > uuid2[i] {
			t.Errorf("UUID v7 ordering violation: later UUID has earlier timestamp")
			break
		}
		// If bytes are equal, continue to next byte
	}
	
	t.Logf("UUID V7 (earlier): %s", uuid1)
	t.Logf("UUID V7 (later):   %s", uuid2)
}

func TestV7TimestampExtraction(t *testing.T) {
	beforeTime := time.Now().UnixMilli()
	uuid := NewV7()
	afterTime := time.Now().UnixMilli()
	
	// Extract timestamp from UUID (first 48 bits)
	// Direct extraction from the first 8 bytes, then shift right by 16 bits
	extractedTimestamp := int64(binary.BigEndian.Uint64(uuid[0:8]) >> 16)
	
	if extractedTimestamp < beforeTime || extractedTimestamp > afterTime {
		t.Errorf("extracted timestamp %d not within expected range [%d, %d]", 
			extractedTimestamp, beforeTime, afterTime)
	}
	
	t.Logf("Generated timestamp: %d", extractedTimestamp)
	t.Logf("Expected range: [%d, %d]", beforeTime, afterTime)
}

func BenchmarkNewV7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV7()
	}
}