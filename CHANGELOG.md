# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **UUID v7 support** - Implementation of UUID version 7 according to RFC 9562
  - `NewV7()` function for generating timestamp-based UUIDs with millisecond precision
  - 48-bit timestamp (milliseconds since Unix epoch) with random data for uniqueness
  - Monotonic ordering by creation time
  - High performance generation (126.3 ns/op)
  - Comprehensive test suite including uniqueness, ordering, and timestamp extraction tests

### Changed
- Updated README.md with UUID v7 documentation and usage examples
- Updated benchmarks to include UUID v7 performance metrics
- Added RFC 9562 reference for UUID v7 specification

### Technical Details
- UUID v7 structure: 48-bit timestamp + 12-bit random + version + variant + 62-bit random
- Maintains compatibility with existing UUID v1-v5 implementations
- All tests passing with no breaking changes to existing functionality