# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.17.0] - 2026-03-01

### Added
- Geo types support via `github.com/twpayne/go-geom` (#404)
- Semver type with proper query filters (#457)
- Custom field name override via Go struct tag (#456)
- Proper struct sub-filter queries (#453)
- Insert method for bulk record creation (#452)
- Index rebuild method for repos (#451)
- Fulltext search any (or) match (#449)
- Support for `time.Month` and `time.Weekday` types (#448)
- Count index added to all table schemas (#445)

### Fixed
- Go tags syntax alignment (#455)
- Datetime CBOR handling (#450)
- Query rename offset alignment (#447)
- Remove live query variable support workaround (#446)

### Changed
- Require surrealdb.go client version 1.3.0 (#458)

### Documentation
- Comprehensive documentation overhaul (#460)

## [v0.16.0] - 2026-02-23

### Added
- Support for complex ID types (array and object IDs) (#430)
- Complex ID range queries (#443)

### Changed
- **Breaking:** Overhaul of parser setup (#441)
- Optimised CLI command and flags (#442)
- Internal refactor of `subPkg` to `relativePkgPath` (#440)

## [v0.15.0] - 2026-02-21

### Added
- SurrealDB v3.0 compatibility (#438)
- String record ID with configurable generator (#428)

### Changed
- Optimised table ID handling (#427)
- Optimised Go template handling (#433)

## [v0.14.0] - 2026-01-30

### Added
- Lifecycle hooks for repo and model (#423)
- Query fetch for distinct field values (#422)
- **Breaking:** Query builder wording aligned with official SurrealQL (#418)

### Fixed
- Slices with nil values (#425)
- Nil pointer dereference in repo read for missing records (#421)
- Errors not exposed properly (#419)

## [v0.13.0] - 2026-01-26

### Added
- Optional soft-delete handling for nodes (#410)
- Opt-in context cache for database reads (#409)
- Wire dependency injection setup generation (#417)
- Fetch support for live queries (#406)
- Support for `github.com/gofrs/uuid` (#403)
- Most missing field filter functions (#402)

## [v0.12.0] - 2025-11-30

### Added
- Indexing, unique constraints and fulltext search (#153)
- Fulltext search features for query builder (#397)
- Tempfiles statement for query builder (#399)

## [v0.11.0] - 2025-11-26

### Added
- Optimistic locking feature for models (#99)
- Smart password field type (#300)
- Email field type (#390)
- Iterator methods for query builder output (#389)
- Most field types are now sortable (#387)
- Code coverage analysis (#219)
- Initial GitBook documentation (#305)

### Changed
- `created_at` fields now use `readonly` in schema (#385)

## [v0.10.0] - 2025-11-23

### Added
- Empty slice filter, schema and marshal support (#373)

### Changed
- **Breaking:** Better conversion layer with direct CBOR marshal/unmarshal (#374)
- Cleanup of legacy conversion methods (#381)
- Switched from asdf to mise for version management (#372)

### Fixed
- Schema generation for deeply nested fields (#382)
- Live query race condition (#377)

## [v0.9.0] - 2025-11-18

### Changed
- **Breaking:** Switch driver from sdbc to official surrealdb.go client (#371)

## [v0.8.0] - 2025-05-21

### Changed
- **Breaking:** SOM is now a compile-time only dependency (#359)
- **Breaking:** CLI command moved to root package (#361)
- Upgraded `github.com/urfave/cli` to v3 (#354)

### Fixed
- Config type alias to sdbc (#358)
- `go.mod` commands failing with embedded template files (#357)

## [v0.7.1] - 2025-05-04

### Added
- DB 2.0 enforced relations, enum string literals and more functions (#328)

### Changed
- Required Go version updated to v1.23.x (#350)
- Required sdbc version updated to v0.9.2 (#349)

## [v0.7.0] - 2024-12-11

### Added
- Support for SurrealDB 2.0 (#321)

### Changed
- **Breaking:** Updated for SurrealDB 2.0 compatibility (#321)
- Required sdbc v0.9.0 (#326)

## [v0.6.4] - 2024-11-30

### Changed
- Updated to sdbc v0.9.0 (#306)

## [v0.6.3] - 2024-09-05

### Fixed
- Codegen for pointer and non-pointer node/edge types (#316)

## [v0.6.2] - 2024-08-28

### Changed
- **Breaking:** `sdbc.ID` renamed to `som.ID` (#311)
- Tests moved to internal package (#310)

## [v0.6.1] - 2024-08-28

### Fixed
- Generated files not written correctly (#309)

## [v0.6.0] - 2024-08-28

### Added
- Field-to-field comparison filters (#302)
- `DescribeWithVars` and `Debug` methods for query builder (#304)
- Missing filter functions (#303)

### Changed
- Static param names in generated code (#301)

## [v0.5.0] - 2024-08-14

### Added
- Go, SOM and sdbc version checks before code generation (#266)
- `time.Duration` type support (#164)
- Filter functions and optimisation (#291)
- Table types in generated schema definition (#264)
- In-memory file system for codegen (#293)

### Changed
- **Breaking:** Switched to CBOR protocol via sdbc update (#289)
- Optimised schema definition for timestamps (#265)
- Optimised CLI setup (#221)

### Fixed
- Hardcoded import path in embedded file (#263)
- Embedded file using local import (#296)

## [v0.4.0] - 2024-05-02

### Added
- Live count method (#234)
- Generated comments for repo methods (#247)

### Fixed
- Filter for live queries (#259)

## [v0.3.0] - 2023-11-27

### Added
- Refresh method for repos (#237)
- Support for `byte` and `[]byte` fields (#119)
- Support for `net/url.URL` type (#163)
- Support for missing native numeric types (#174)

### Changed
- Reduced codegen by making query builder generic (#245)
- Moved generic CRUD methods to embed (#246)

## [v0.2.0] - 2023-11-13

### Fixed
- Generated filters missing methods for edges (#243)
- Generated names for types and fields (#238)
- Duplicates in generated code (#240)
- Edge field conversion (#241)
- Nil dereference due to invalid use of anonymous struct field type (#239)

## [v0.1.2] - 2023-10-16

### Fixed
- Timestamps zero value dereference nil pointer (#227)

## [v0.1.1] - 2023-10-09

### Added
- Client options passthrough to sdbc (#222)

## [v0.1.0] - 2023-09-17

### Added
- ULID as primary key option for nodes (#204)
- License file (#208)
- Security policy (#172)

### Changed
- **Breaking:** Repository moved to `go-surreal` org (#194)
- **Breaking:** sdbc (SurrealDB client) extracted to own repo (#203)
- Timestamp handling moved to database schema layer (#197)
- Tested against SurrealDB v1.0.0 (#207)

## [v0.0.11] - 2023-09-05

### Added
- Live queries support (#167)
- SurrealDB strict-mode compatibility (#175)
- Custom SurrealDB client implementation (#152)
- Support for SurrealDB beta.11 (#200)

### Fixed
- Enum and struct fields handling (#199)
- Schema with transaction broken (#186)

## [v0.0.10] - 2023-08-04

### Added
- Async methods for query builder (#171)
- Transaction wrapping for schema definition (#156)
- ID field definition in generated schema (#165)
- UUID field assertion in schema (#162)

### Changed
- Moved static code to embed files (#170)

### Fixed
- Schema needs `option<?>` types for nilable fields (#161)
- Count query and result mapping (#130)
- Unexpected DB response structure for unmarshal (#131)
- Query variable index breaks after 26 variables (#75)

### Performance
- Optimised query builder string concatenation (#173)

## [v0.0.9] - 2023-04-09

### Changed
- Updated for SurrealDB beta.9 (#129)
- Updated surrealdb.go client to v0.2.0 (#125)

## [v0.0.8] - 2023-03-29

### Added
- Generated code is now mockable via interfaces (#112)
- Lib package moved into generated sources via embed (#118)

## [v0.0.7] - 2023-02-13

### Fixed
- Record links (#98)

## [v0.0.6] - 2023-01-25

### Added
- Sub-queries for nodes and edges (#84)

## [v0.0.5] - 2023-01-19

### Added
- Schema asserts for enum values (#83)
- Asserts for database schema definition (#82)

## [v0.0.4] - 2022-12-29

### Added
- Auto timestamps on create/update (#79)
- Query describe method (#74)
- Reusable query builder (#76)

### Changed
- **Breaking:** Separate `CreateWithID` method provided (#80)

## [v0.0.3] - 2022-12-26

### Added
- Database schema generation for strict mode (#44)
- Select, update and delete operations (#57)
- Pointer field support (#64)
- Custom record IDs (#68)
- "Do not edit" comment in generated code (#63)
- Testcontainers integration testing (#62)

## [v0.0.2] - 2022-12-08

### Added
- CI workflow for pull requests (#43)
- Conventional commits enforcement (#37)
- Golangci-lint integration (#51)

### Fixed
- Lint issues in example/gen (#56)

## [v0.0.1] - 2022-12-03

Initial release.

### Added
- Basic code generation from Go struct models
- Parser for `som.Node` and `som.Edge` types
- Query builder with fetch statement
- Edge (graph) connection support
- CLI tool for code generation

[Unreleased]: https://github.com/go-surreal/som/compare/v0.17.0...HEAD
[v0.17.0]: https://github.com/go-surreal/som/compare/v0.16.0...v0.17.0
[v0.16.0]: https://github.com/go-surreal/som/compare/v0.15.0...v0.16.0
[v0.15.0]: https://github.com/go-surreal/som/compare/v0.14.0...v0.15.0
[v0.14.0]: https://github.com/go-surreal/som/compare/v0.13.0...v0.14.0
[v0.13.0]: https://github.com/go-surreal/som/compare/v0.12.0...v0.13.0
[v0.12.0]: https://github.com/go-surreal/som/compare/v0.11.0...v0.12.0
[v0.11.0]: https://github.com/go-surreal/som/compare/v0.10.0...v0.11.0
[v0.10.0]: https://github.com/go-surreal/som/compare/v0.9.0...v0.10.0
[v0.9.0]: https://github.com/go-surreal/som/compare/v0.8.0...v0.9.0
[v0.8.0]: https://github.com/go-surreal/som/compare/v0.7.1...v0.8.0
[v0.7.1]: https://github.com/go-surreal/som/compare/v0.7.0...v0.7.1
[v0.7.0]: https://github.com/go-surreal/som/compare/v0.6.4...v0.7.0
[v0.6.4]: https://github.com/go-surreal/som/compare/v0.6.3...v0.6.4
[v0.6.3]: https://github.com/go-surreal/som/compare/v0.6.2...v0.6.3
[v0.6.2]: https://github.com/go-surreal/som/compare/v0.6.1...v0.6.2
[v0.6.1]: https://github.com/go-surreal/som/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/go-surreal/som/compare/v0.5.0...v0.6.0
[v0.5.0]: https://github.com/go-surreal/som/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/go-surreal/som/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/go-surreal/som/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/go-surreal/som/compare/v0.1.2...v0.2.0
[v0.1.2]: https://github.com/go-surreal/som/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/go-surreal/som/compare/v0.1.0...v0.1.1
[v0.1.0]: https://github.com/go-surreal/som/compare/v0.0.11...v0.1.0
[v0.0.11]: https://github.com/go-surreal/som/compare/v0.0.10...v0.0.11
[v0.0.10]: https://github.com/go-surreal/som/compare/v0.0.9...v0.0.10
[v0.0.9]: https://github.com/go-surreal/som/compare/v0.0.8...v0.0.9
[v0.0.8]: https://github.com/go-surreal/som/compare/v0.0.7...v0.0.8
[v0.0.7]: https://github.com/go-surreal/som/compare/v0.0.6...v0.0.7
[v0.0.6]: https://github.com/go-surreal/som/compare/v0.0.5...v0.0.6
[v0.0.5]: https://github.com/go-surreal/som/compare/v0.0.4...v0.0.5
[v0.0.4]: https://github.com/go-surreal/som/compare/v0.0.3...v0.0.4
[v0.0.3]: https://github.com/go-surreal/som/compare/v0.0.2...v0.0.3
[v0.0.2]: https://github.com/go-surreal/som/compare/v0.0.1...v0.0.2
[v0.0.1]: https://github.com/go-surreal/som/releases/tag/v0.0.1
