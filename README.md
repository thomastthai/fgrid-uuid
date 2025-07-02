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
  }
  ```

## benchmarks
  ```
  BenchmarkNewV1	 2743015	       437.4 ns/op
  BenchmarkNewV3	 5820739	       205.4 ns/op
  BenchmarkNewV4	12387554	        95.81 ns/op
  BenchmarkNewV5	 4918503	       243.3 ns/op
  BenchmarkNewV7	 9431486	       126.3 ns/op
  ```

## UUID versions

### Version 1 (timestamp + MAC)
Based on timestamp and MAC address. Provides uniqueness across space and time.

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

## documentation
* @[Sourcegraph](http://sourcegraph.com/github.com/thomastthai/fgrid-uuid)

## links
* [RFC 4122](http://tools.ietf.org/html/rfc4122) - Original UUID specification
* [RFC 9562](https://tools.ietf.org/html/rfc9562) - Updated UUID specification with version 7

## badges
[![status](https://sourcegraph.com/api/repos/github.com/fgrid/uuid/.badges/status.svg)](https://sourcegraph.com/github.com/fgrid/uuid) [![library users](https://sourcegraph.com/api/repos/github.com/fgrid/uuid/.badges/library-users.svg)](https://sourcegraph.com/github.com/fgrid/uuid) [![dependents](https://sourcegraph.com/api/repos/github.com/fgrid/uuid/.badges/dependents.svg)](https://sourcegraph.com/github.com/fgrid/uuid) [![views](https://sourcegraph.com/api/repos/github.com/fgrid/uuid/.counters/views.svg)](https://sourcegraph.com/github.com/fgrid/uuid) [![Go Report Card](http://goreportcard.com/badge/fgrid/uuid)](http://goreportcard.com/report/fgrid/uuid)
