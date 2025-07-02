[![Stories in Ready](https://badge.waffle.io/fgrid/uuid.png?label=ready&title=Ready)](https://waffle.io/fgrid/uuid)
# uuid
golang uuid generator

Supports UUID versions 1, 3, 4, 5, and 7 according to RFC 4122 and RFC 9562.

## install
  ```
  go get github.com/thomastthai/fgrid-uuid
  ```

## usage
  ```go
  package main

  import (
      "fmt"
      "encoding/json"
      "github.com/thomastthai/fgrid-uuid"
  )

  func main() {
      // UUID v1 - timestamp and MAC address based
      v1 := uuid.NewV1()
      fmt.Printf("UUID v1: %s\n", v1.String())

      // UUID v3 - namespace and name with MD5 hash
      v3 := uuid.NewV3(uuid.NameSpaceDNS, []byte("example.com"))
      fmt.Printf("UUID v3: %s\n", v3.String())

      // UUID v4 - random
      v4 := uuid.NewV4()
      fmt.Printf("UUID v4: %s\n", v4.String())

      // UUID v5 - namespace and name with SHA-1 hash
      v5 := uuid.NewV5(uuid.NameSpaceDNS, []byte("example.com"))
      fmt.Printf("UUID v5: %s\n", v5.String())

      // UUID v7 - timestamp with random data (RFC 9562)
      v7 := uuid.NewV7()
      fmt.Printf("UUID v7: %s\n", v7.String())
      
      // Parse UUID from string
      parsed, err := uuid.ParseUUID("550e8400-e29b-41d4-a716-446655440000")
      if err != nil {
          panic(err)
      }
      fmt.Printf("Parsed UUID: %s\n", parsed.String())
      
      // JSON marshaling (implements encoding.TextMarshaler)
      data, _ := json.Marshal(v4)
      fmt.Printf("JSON: %s\n", data)
      
      // JSON unmarshaling (implements encoding.TextUnmarshaler)
      var unmarshaled uuid.UUID
      json.Unmarshal(data, &unmarshaled)
      fmt.Printf("Unmarshaled: %s\n", unmarshaled.String())
  }
  ```

### Advanced Usage

#### Custom V1 Generator
For applications that need explicit control over UUID v1 generation lifecycle:

  ```go
  // Create a custom V1 generator
  generator := uuid.NewUUIDV1Generator()
  
  // Generate UUIDs using the custom generator
  uuid1 := generator.NewV1()
  uuid2 := generator.NewV1()
  
  // Multiple generators can coexist
  generator2 := uuid.NewUUIDV1Generator()
  uuid3 := generator2.NewV1()
  ```

## benchmarks
  ```
  BenchmarkNewV1                       	 2679056	       445.3 ns/op
  BenchmarkUUIDV1Generator_NewV1       	 2740653	       432.2 ns/op
  BenchmarkUUIDV1Generator_Concurrent  	 2822744	       427.8 ns/op
  BenchmarkNewV3                       	 5772448	       206.9 ns/op
  BenchmarkNewV4                       	13480060	        88.19 ns/op
  BenchmarkNewV5                       	 4902060	       246.4 ns/op
  BenchmarkNewV7                       	 9484611	       125.0 ns/op
  ```

## UUID versions

### Version 1 (timestamp + MAC)
Based on timestamp and MAC address. Provides uniqueness across space and time.

**Thread Safety**: All UUID generation functions are thread-safe and can be called concurrently from multiple goroutines. The V1 generator uses a background goroutine to ensure timestamp monotonicity and proper clock sequence handling according to RFC 4122.

**Custom Generators**: For advanced use cases, you can create custom V1 generators using `NewUUIDV1Generator()`. Each generator maintains its own state and can be used independently.

### Version 3 (namespace + name + MD5)
Based on namespace UUID and name with MD5 hash. Deterministic.

### Version 4 (random)
Based on random or pseudo-random numbers. Most commonly used.

### Version 5 (namespace + name + SHA-1)
Based on namespace UUID and name with SHA-1 hash. Deterministic.

### Version 7 (timestamp + random)
**NEW**: Based on Unix timestamp in milliseconds with random data. Provides:
- Monotonic ordering by creation time
- High performance generation
- 48-bit timestamp (good until year 10,895 CE)
- RFC 9562 compliance

## Features

### Go Interface Support
- **fmt.Stringer**: UUIDs can be printed directly with `fmt.Printf`
- **encoding.TextMarshaler**: JSON marshaling support
- **encoding.TextUnmarshaler**: JSON unmarshaling support
- **ParseUUID**: Flexible parsing from string (supports both hyphenated and non-hyphenated formats)

### Thread Safety
All functions are designed for concurrent use:
- UUID generation functions are thread-safe
- V1 generator maintains proper timestamp ordering even under high concurrency
- Multiple V1 generators can coexist without interference

## documentation
* @[Sourcegraph](http://sourcegraph.com/github.com/thomastthai/fgrid-uuid)

## links
* [RFC 4122](http://tools.ietf.org/html/rfc4122) - Original UUID specification
* [RFC 9562](https://tools.ietf.org/html/rfc9562) - Updated UUID specification with version 7

## badges
[![status](https://sourcegraph.com/api/repos/github.com/fgrid/uuid/.badges/status.svg)](https://sourcegraph.com/github.com/fgrid/uuid) [![library users](https://sourcegraph.com/api/repos/github.com/fgrid/uuid/.badges/library-users.svg)](https://sourcegraph.com/github.com/fgrid/uuid) [![dependents](https://sourcegraph.com/api/repos/github.com/fgrid/uuid/.badges/dependents.svg)](https://sourcegraph.com/github.com/fgrid/uuid) [![views](https://sourcegraph.com/api/repos/github.com/fgrid/uuid/.counters/views.svg)](https://sourcegraph.com/github.com/fgrid/uuid) [![Go Report Card](http://goreportcard.com/badge/fgrid/uuid)](http://goreportcard.com/report/fgrid/uuid)
