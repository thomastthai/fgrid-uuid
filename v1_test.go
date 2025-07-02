package uuid

import (
	"sync"
	"testing"
	"time"
)

func TestNewV1(t *testing.T) {
	uuid := NewV1()
	if uuid.Version() != 1 {
		t.Errorf("invalid version %d - expected 1", uuid.Version())
	}
	t.Logf("UUID V1: %s", uuid)
}

func BenchmarkNewV1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV1()
	}
}

func TestNewUUIDV1Generator(t *testing.T) {
	gen := NewUUIDV1Generator()
	if gen == nil {
		t.Fatal("NewUUIDV1Generator() returned nil")
	}
	
	if len(gen.mac) != 6 {
		t.Errorf("MAC address length = %d, want 6", len(gen.mac))
	}
	
	if !gen.started {
		t.Error("Generator not marked as started")
	}
}

func TestUUIDV1Generator_NewV1(t *testing.T) {
	gen := NewUUIDV1Generator()
	uuid := gen.NewV1()
	
	if uuid == nil {
		t.Fatal("NewV1() returned nil")
	}
	
	if uuid.Version() != 1 {
		t.Errorf("invalid version %d - expected 1", uuid.Version())
	}
	
	t.Logf("Generated UUID V1: %s", uuid)
}

func TestUUIDV1Generator_Concurrent(t *testing.T) {
	gen := NewUUIDV1Generator()
	const numGoroutines = 100
	const uuidsPerGoroutine = 10
	
	var wg sync.WaitGroup
	results := make(chan *UUID, numGoroutines*uuidsPerGoroutine)
	
	// Start multiple goroutines generating UUIDs concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < uuidsPerGoroutine; j++ {
				uuid := gen.NewV1()
				results <- uuid
			}
		}()
	}
	
	wg.Wait()
	close(results)
	
	// Collect all UUIDs and check for uniqueness
	uuids := make(map[string]bool)
	count := 0
	for uuid := range results {
		if uuid == nil {
			t.Error("Generated nil UUID")
			continue
		}
		
		if uuid.Version() != 1 {
			t.Errorf("invalid version %d - expected 1", uuid.Version())
		}
		
		str := uuid.String()
		if uuids[str] {
			t.Errorf("duplicate UUID generated: %s", str)
		}
		uuids[str] = true
		count++
	}
	
	expectedCount := numGoroutines * uuidsPerGoroutine
	if count != expectedCount {
		t.Errorf("generated %d UUIDs, expected %d", count, expectedCount)
	}
	
	t.Logf("Successfully generated %d unique UUIDs concurrently", count)
}

func TestUUIDV1Generator_MultipleGenerators(t *testing.T) {
	// Test that multiple generators can coexist
	gen1 := NewUUIDV1Generator()
	gen2 := NewUUIDV1Generator()
	
	uuid1 := gen1.NewV1()
	uuid2 := gen2.NewV1()
	
	if uuid1.String() == uuid2.String() {
		t.Error("Different generators produced identical UUIDs")
	}
	
	// Both should be version 1
	if uuid1.Version() != 1 || uuid2.Version() != 1 {
		t.Error("Generated UUIDs are not version 1")
	}
	
	t.Logf("Generator 1 UUID: %s", uuid1)
	t.Logf("Generator 2 UUID: %s", uuid2)
}

func TestUUIDV1Generator_TimestampMonotonicity(t *testing.T) {
	gen := NewUUIDV1Generator()
	
	// Generate UUIDs in quick succession
	var timestamps []uint64
	for i := 0; i < 10; i++ {
		uuid := gen.NewV1()
		// Extract timestamp from UUID (first 8 bytes, but we need to reconstruct the original format)
		timestamp := uint64(uuid[6])<<8 | uint64(uuid[7])          // low bits
		timestamp |= uint64(uuid[4])<<24 | uint64(uuid[5])<<16     // mid bits  
		timestamp |= uint64(uuid[0])<<56 | uint64(uuid[1])<<48 | uint64(uuid[2])<<40 | uint64(uuid[3])<<32 // high bits
		timestamp &= 0x0fffffffffffffff // Clear version bits
		timestamps = append(timestamps, timestamp)
		
		// Small delay to potentially get different timestamps
		time.Sleep(time.Microsecond)
	}
	
	// Check that timestamps are non-decreasing (monotonic)
	for i := 1; i < len(timestamps); i++ {
		if timestamps[i] < timestamps[i-1] {
			t.Errorf("timestamp decreased: %d -> %d at position %d", timestamps[i-1], timestamps[i], i)
		}
	}
	
	t.Logf("Generated %d UUIDs with monotonic timestamps", len(timestamps))
}

func BenchmarkUUIDV1Generator_NewV1(b *testing.B) {
	gen := NewUUIDV1Generator()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		gen.NewV1()
	}
}

func BenchmarkUUIDV1Generator_Concurrent(b *testing.B) {
	gen := NewUUIDV1Generator()
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gen.NewV1()
		}
	})
}
