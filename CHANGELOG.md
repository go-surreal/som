# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.19.0] - 2026-06-21


### Added

- Require surrealdb v3.0.3 (#477)
- Require surrealdb v3.0.5 (#483)
- Optimise live query handling (#478)
- Not filter for query builder (#485)
- Run against surrealdb v3.1.0 (#499)
- Run against surrealdb v3.1.3 (#502)
- Run against surrealdb v3.1.5 (#505)

### Fixed

- Circumvent surrealdb.go live query race condition (#488)
- Bump surrealdb.go and drop live query race workaround (#495)
- Remove unnecessary som version checks (#496)

## [v0.18.0] - 2026-03-09


### Added

- Add support for client-side transactions (#463)
- Surrealdb v3 structured server errors (#465)
- Raw queries (#466)
- Range query for string id nodes (#471)
- Where field value between (#472)
- Verify db version when connecting (#474)

### Changed

- Remove db client wrapper & cbor round-trip (#468)

### Fixed

- Surrealdb v3.0.2 query search same field  (#462)

### Documentation

- Prepare release v0.18.0 (#475)

## [v0.17.0] - 2026-03-01


### Added

- Add count index to all tables (#445)
- Fulltext search any (or) match (#449)
- Add index rebuild to repo methods (#451)
- Add insert method to repo for bulk create (#452)
- Add support for time.Month and time.Weekday (#448)
- Support proper struct sub-filter queries (#453)
- Add semver type with proper query filters (#457)
- Allow explicit custom field names via go tag (#456)
- Add geo types support (#404)
- Require surrealdb.go client version 1.3.0 (#458)

### Fixed

- Remove live query variable support workaround (#446)
- Align query rename offset to start (#447)
- Datetime cbor handling (#450)
- Properly align go tags syntax (#455)

### Documentation

- Update deprecated and add missing doc files (#460)

## [v0.16.0] - 2026-02-23


### Added

- Add support for complex id types (#430)
- Optimise cli command and flags (#442)
- Complex id range queries  (#443)

### Changed

- SubPkg to relativePkgPath (#440)
- Overhaul parser setup (#441)

## [v0.15.0] - 2026-02-21


### Added

- Optimise table id handling (#427)
- String record id with configurable generator (#428)
- Surrealdb v3.0 compatibility (#438)

### Changed

- Code cleanup & more tests (#426)
- Optimise go template handling (#433)

## [v0.14.0] - 2026-01-30


### Added

- Align query builder wording with official query lang (#418)
- Query fetch distinct field values (#422)
- Add lifecycle hooks to repo and model (#423)

### Fixed

- Expose errors properly (#419)
- Prevent nil pointer deref in repo read for missing records (#421)
- Slices with nil values (#425)

## [v0.13.0] - 2026-01-26


### Added

- Add most missing field filter functions (#402)
- Add support for github.com/gofrs/uuid (#403)
- Allow fetch for live queries (#406)
- Add opt-in context cache for database reads (#409)
- Optional soft-delete handling for nodes (#410)
- Generate wire setup if used by project (#417)

### Changed

- Cleanup query builder type (#416)

## [v0.12.0] - 2025-11-30


### Added

- Add indexing, unique and fulltext (#153)
- Add fulltext search features to query builder (#397)
- Add tempfiles statement to query builder (#399)

### Documentation

- Add fulltext search features (#398)

## [v0.11.0] - 2025-11-26


### Added

- Use readonly for created_at fields (#385)
- Make most field types sortable (#387)
- Add email field type (#390)
- Add iterator methods to query builder output (#389)
- Add smart password field type (#300)
- Optimistic locking feature for models (#99)

### Changed

- Remove unused fieldDef methods from fields (#393)

### Fixed

- Navigation and summary (#391)

### Documentation

- Initial documentation setup (#305)
- Add new features to gitbook (#394)

## [v0.10.0] - 2025-11-23


### Added

- Empty filter, schema and marshal fixup, testing (#373)

### Changed

- Better conv layer with direct cbor (un)marshal (#374)
- Cleanup legacy conv methods (#381)

### Fixed

- Schema generation for deeply nested fields (#382)

## [v0.9.0] - 2025-11-18


### Added

- Switch driver from sdbc to official surrealdb go client (#371)

### Fixed

- Main workflow wrong working dir (#362)

## [v0.8.0] - 2025-05-21


### Added

- Make som a compile-time only dependency (#359)
- Move cmd to root (#361)

### Fixed

- Go.mod commands fail with embedded template files (#357)
- Make config type alias to sdbc (#358)

## [v0.7.1] - 2025-05-04


### Added

- Db 2.0 enforced relations, enum string literals & more funcs (#328)

## [v0.7.0] - 2024-12-11


### Added

- Support for surrealdb 2.0 (#321)

## [v0.6.4] - 2024-11-30


### Fixed

- Testcontainers terminate does not ignore unnecessary error (#323)

## [v0.6.3] - 2024-09-05


### Fixed

- Codegen for (no-)pointer node/edge types (#316)

## [v0.6.2] - 2024-08-28


### Changed

- Hide tests in internal package (#310)
- Sdbc.ID to som.ID (#311)

## [v0.6.1] - 2024-08-27


### Fixed

- Generated files not written correctly (#309)

## [v0.6.0] - 2024-08-27


### Added

- Filter by direct field to field comparisons (#302)
- Add DescribeWithVars and Debug methods to query builder (#304)
- Add missing filter functions (#303)

### Changed

- Static param name in generated code (#301)

## [v0.5.0] - 2024-08-14


### Added

- Add table types to generated schema definition (#264)
- Optimise schema definition for timestamps (#265)
- Update sdbc & switch to cbor protocol (#289)
- Support time.duration type (#164)
- Implement filter functions and optimisation (#291)
- In-memory file system for codegen (#293)
- Check go, som and sdbc version before generate (#266)

### Changed

- Optimise cli setup (#221)
- Code cleanup (#292)
- Embedded gen files handling (#298)

### Fixed

- Hardcoded import path in embedded file (#263)
- Embedded file uses local import (#296)

## [v0.4.0] - 2024-05-02


### Added

- Implement live count method (#234)
- Generate comments for repo methods (#247)

### Fixed

- Filter for live queries (#259)

## [v0.3.0] - 2023-11-27


### Added

- Add support for byte and []byte fields (#119)
- Add support for net/url.URL type (#163)
- Add support for missing native numeric types (#174)
- Implement refresh method for repos (#237)

### Changed

- Reduce codegen by making query builder generic (#245)
- Move generic crud methods to embed (#246)

## [v0.2.0] - 2023-11-13


### Fixed

- Dependabot bump google.golang.org/grpc to 1.57.1 (#231)
- Dependabot bump github.com/docker/docker to 24.0.7 (#233)
- Nil deref due to invalid use of anonymous struct field type (#239)
- Edge field conversion (#241)
- Remove duplicates from generated code (#240)
- Generated names for types & fields (#238)
- Generated filters missing methods for edges (#243)

## [v0.1.2] - 2023-10-16


### Fixed

- Timestamps zero value deref nil pointer (#227)

## [v0.1.1] - 2023-10-09


### Added

- Add client options and pass to sdbc (#222)

## [v0.1.0] - 2023-09-17


### Added

- Create nodes with ulid as primary key (#204)

### Changed

- Move repo to go-surreal org (#194)
- Move sdbc code to own repo (#203)
- Move timestamp handling to database schema layer (#197)

### Fixed

- Codegen bugs & streamline tests (#210)

### Documentation

- Update readme (#201)
- Update readme & add ideas (#160)

## [v0.0.11] - 2023-09-05


### Added

- Add surrealdb strict-mode compatibility (#175)
- Implement custom surrealdb client (#152)
- Optimise sdbc close & cleanup handling (#191)
- Add support for surrealdb beta.11 (#200)
- Add support for live queries (#167)

### Changed

- Let conv functions use pointers (#198)

### Fixed

- Apply schema with transaction is broken (#186)
- Apply latest surrealdb nightly changes (#187)
- Enum & struct fields handling (#199)

## [v0.0.10] - 2023-08-04


### Added

- Add transaction wrapping to schema definition (#156)
- Add assert for uuid fields to schema (#162)
- Add id field definition to generated schema (#165)
- Add async methods to query builder (#171)

### Changed

- Simplify conv of uuid type (#166)
- Move static code to embed files (#170)
- Optimise query builder string concat (#173)

### Fixed

- Count query and result mapping (#130)
- Unexpected db response structure for unmarshal (#131)
- Dependabot alert by bumping google.golang.org/grpc to 1.53.0 (#145)
- Tests (#157)
- Schema needs to set option<?> types for nilable fields (#161)

## [v0.0.9] - 2023-04-09


### Added

- Utilize new features of surrealdb.go client update v0.2.0 (#125)
- Update for surrealdb version beta.9 (#129)

### Fixed

- Dependabot alert by updating opencontainers/runc to 1.1.5 (#123)

## [v0.0.8] - 2023-03-29


### Added

- Make generated code mockable by using interfaces (#112)
- Small gen updates, more tests and examples (#106)
- Move lib package into generated sources via embed (#118)

### Documentation

- Update readme & faq (#124)

## [v0.0.7] - 2023-02-13


### Fixed

- Record links (#98)

### Documentation

- Add ideas document (#92)

## [v0.0.6] - 2023-01-24


### Added

- Implement sub-queries for nodes and edges (#84)

## [v0.0.5] - 2023-01-19


### Added

- Add asserts to database schema definition (#82)
- Schema assert for enum values (#83)

### Documentation

- Update readme (#72)

## [v0.0.4] - 2022-12-29


### Added

- Implement query describe method (#74)
- Make query builder reusable (#76)
- Implement auto timestamps on create/update (#79)

### Changed

- Provide separate create-with-id method (#80)

### Fixed

- Query variable index breaks after 26 variables (#75)

## [v0.0.3] - 2022-12-26


### Added

- Implement select, update and delete (#57)
- Add "do not edit" comment to generated code (#63)
- Allow custom record ids (#68)
- Support pointer fields (#64)
- Generate database schema for strict usage (#44)

### Changed

- Update database interface (#58)
- Code generator for better extensibility (#66)
- Code generation streamlining (#67)

## [v0.0.2] - 2022-12-08


### Added

- Enforce conventional commits spec via PR workflow (#37)

### Fixed

- Golangci-lint issues (#52)
- Lint example/gen and handle existing issues (#56)

### Documentation

- Update README (#55)

## [v0.0.1] - 2022-12-03

[v0.19.0]: https://github.com/go-surreal/som/compare/v0.18.0...v0.19.0
[v0.18.0]: https://github.com/go-surreal/som/compare/v0.17.0...v0.18.0
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

