package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"net"
	"sync"
	"time"
)

type stamp [10]byte

// UUIDV1Generator encapsulates the state and logic for generating UUID version 1.
// It maintains the MAC address, timestamp state, and clock sequence in a thread-safe manner.
// The generator uses a background goroutine to ensure timestamp monotonicity and proper
// clock sequence handling as required by RFC 4122.
type UUIDV1Generator struct {
	mac      []byte
	requests chan bool
	answers  chan stamp
	once     sync.Once
	started  bool
}

const gregorianUnix = 122192928000000000 // nanoseconds between gregorion zero and unix zero

// Package-level singleton for backward compatibility
var (
	defaultV1Generator     *UUIDV1Generator
	defaultV1GeneratorOnce sync.Once
)

// getDefaultV1Generator returns the singleton V1 generator, creating it if necessary.
// This ensures backward compatibility with the existing NewV1() function.
func getDefaultV1Generator() *UUIDV1Generator {
	defaultV1GeneratorOnce.Do(func() {
		defaultV1Generator = NewUUIDV1Generator()
	})
	return defaultV1Generator
}

// NewUUIDV1Generator creates a new UUID version 1 generator with proper initialization.
// It discovers the MAC address from network interfaces, starts the background timestamp
// goroutine, and returns a generator ready for concurrent use.
// 
// The generator is thread-safe and can be used concurrently from multiple goroutines.
// Each generator maintains its own clock sequence and timestamp state to ensure
// uniqueness even when system time moves backward.
func NewUUIDV1Generator() *UUIDV1Generator {
	g := &UUIDV1Generator{
		mac:      make([]byte, 6),
		requests: make(chan bool),
		answers:  make(chan stamp),
	}
	
	// Initialize with random MAC if no interface found
	rand.Read(g.mac)
	
	// Try to find a real MAC address
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) == 6 {
				copy(g.mac, iface.HardwareAddr[:6])
				break
			}
		}
	}
	
	// Start the background goroutine
	go g.unique()
	g.started = true
	
	return g
}

// NewV1 generates a new UUID version 1 using this generator.
// This method is thread-safe and can be called concurrently.
func (g *UUIDV1Generator) NewV1() *UUID {
	var uuid UUID
	g.requests <- true
	s := <-g.answers
	copy(uuid[:4], s[4:])
	copy(uuid[4:6], s[2:4])
	copy(uuid[6:8], s[:2])
	uuid[6] = (uuid[6] & 0x0f) | 0x10
	copy(uuid[8:10], s[8:])
	copy(uuid[10:], g.mac)
	uuid.variantRFC4122()
	return &uuid
}

// unique runs in a background goroutine to maintain timestamp monotonicity
// and proper clock sequence handling. It ensures that even if the system
// clock moves backward, UUIDs remain unique and properly ordered.
func (g *UUIDV1Generator) unique() {
	var (
		lastNanoTicks uint64
		clockSequence [2]byte
	)
	rand.Read(clockSequence[:])

	for range g.requests {
		var s stamp
		nanoTicks := uint64((time.Now().UTC().UnixNano() / 100) + gregorianUnix)
		if nanoTicks < lastNanoTicks {
			// Clock moved backward, use new random clock sequence
			lastNanoTicks = nanoTicks
			rand.Read(clockSequence[:])
		} else if nanoTicks == lastNanoTicks {
			// Same timestamp, increment to maintain uniqueness
			lastNanoTicks = nanoTicks + 1
		} else {
			// Normal case: time moved forward
			lastNanoTicks = nanoTicks
		}
		binary.BigEndian.PutUint64(s[:], lastNanoTicks)
		copy(s[8:], clockSequence[:])
		g.answers <- s
	}
}

// NewV1 creates a new UUID with variant 1 as described in RFC 4122.
// Variant 1 is based on hosts MAC address and actual timestamp (as count of 100-nanosecond intervals since
// 00:00:00.00, 15 October 1582 (the date of Gregorian reform to the Christian calendar).
// 
// This function maintains backward compatibility by using a singleton generator.
// For explicit control over the generator lifecycle, use NewUUIDV1Generator() instead.
func NewV1() *UUID {
	return getDefaultV1Generator().NewV1()
}
