// Package uuid implements UUID version 1 generation.
// This file contains the implementation for timestamp and MAC address based UUIDs.
package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"net"
	"time"
)

// stamp represents a 10-byte timestamp and clock sequence structure used in UUID v1 generation.
type stamp [10]byte

var (
	// mac stores the MAC address used for UUID v1 generation.
	// If no hardware MAC address is available, a random 6-byte value is used.
	mac []byte
	
	// uuidV1Requests is a channel used to request UUID v1 generation.
	// This ensures thread-safe sequential timestamp generation.
	uuidV1Requests chan bool
	
	// uuidV1Answers is a channel that delivers generated timestamp stamps
	// for UUID v1 creation in a thread-safe manner.
	uuidV1Answers chan stamp
)

// gregorianUnix represents the number of nanoseconds between the Gregorian calendar
// epoch (October 15, 1582) and the Unix epoch (January 1, 1970).
// This constant is used to convert Unix timestamps to UUID v1 timestamps.
const gregorianUnix = 122192928000000000

func init() {
	// Initialize MAC address with random bytes as fallback
	mac = make([]byte, 6)
	rand.Read(mac)
	
	// Initialize communication channels for thread-safe UUID v1 generation
	uuidV1Requests = make(chan bool)
	uuidV1Answers = make(chan stamp)
	
	// Start the goroutine that handles sequential timestamp generation
	go unique()
	
	// Attempt to find a real hardware MAC address
	i, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, d := range i {
		if len(d.HardwareAddr) == 6 {
			mac = d.HardwareAddr[:6]
			return
		}
	}
}

// NewV1 creates a new UUID with version 1 as described in RFC 4122.
// Version 1 UUIDs are based on the host's MAC address and the current timestamp
// (as count of 100-nanosecond intervals since 00:00:00.00, 15 October 1582).
//
// This implementation ensures monotonic timestamps and thread-safety through
// a dedicated goroutine that manages timestamp generation and clock sequence.
func NewV1() *UUID {
	var uuid UUID
	
	// Request a new timestamp stamp from the sequential generator
	uuidV1Requests <- true
	s := <-uuidV1Answers
	
	// Construct UUID according to RFC 4122 Section 4.1.2:
	// - time_low (bytes 0-3): low 32 bits of timestamp
	// - time_mid (bytes 4-5): middle 16 bits of timestamp  
	// - time_hi_and_version (bytes 6-7): high 16 bits of timestamp + version
	// - clock_seq_hi_and_reserved (byte 8): clock sequence + variant
	// - clock_seq_low (byte 9): low 8 bits of clock sequence
	// - node (bytes 10-15): MAC address
	copy(uuid[:4], s[4:])      // time_low
	copy(uuid[4:6], s[2:4])    // time_mid
	copy(uuid[6:8], s[:2])     // time_hi (version will be set below)
	uuid[6] = (uuid[6] & 0x0f) | 0x10  // Set version to 1
	copy(uuid[8:10], s[8:])    // clock_seq
	copy(uuid[10:], mac)       // node (MAC address)
	
	// Set variant bits according to RFC 4122
	uuid.variantRFC4122()
	return &uuid
}

// unique runs in a separate goroutine to ensure thread-safe, monotonic timestamp generation.
// It handles clock sequence management and ensures no duplicate timestamps are generated
// even under high concurrency or clock adjustments.
func unique() {
	var (
		lastNanoTicks uint64    // Last generated timestamp in 100-nanosecond intervals
		clockSequence [2]byte   // Random clock sequence to handle clock adjustments
	)
	
	// Initialize clock sequence with random value
	rand.Read(clockSequence[:])

	// Process timestamp requests sequentially to ensure uniqueness
	for range uuidV1Requests {
		var s stamp
		
		// Convert current time to 100-nanosecond intervals since Gregorian epoch
		nanoTicks := uint64((time.Now().UTC().UnixNano() / 100) + gregorianUnix)
		
		if nanoTicks < lastNanoTicks {
			// Clock moved backwards - update clock sequence and use current time
			lastNanoTicks = nanoTicks
			rand.Read(clockSequence[:])
		} else if nanoTicks == lastNanoTicks {
			// Same timestamp - increment to ensure uniqueness
			lastNanoTicks = nanoTicks + 1
		} else {
			// Normal case - clock moved forward
			lastNanoTicks = nanoTicks
		}
		
		// Pack timestamp into stamp structure (big-endian)
		binary.BigEndian.PutUint64(s[:], lastNanoTicks)
		copy(s[8:], clockSequence[:])
		
		// Send the generated stamp back
		uuidV1Answers <- s
	}
}
