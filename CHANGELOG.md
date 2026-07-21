# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.19.0] - 2026-06-21


### Added

- Require surrealdb v3.0.3 ([#477](https://github.com/go-surreal/som/pull/477))
- Require surrealdb v3.0.5 ([#483](https://github.com/go-surreal/som/pull/483))
- Optimise live query handling ([#478](https://github.com/go-surreal/som/pull/478))
- Not filter for query builder ([#485](https://github.com/go-surreal/som/pull/485))
- Run against surrealdb v3.1.0 ([#499](https://github.com/go-surreal/som/pull/499))
- Run against surrealdb v3.1.3 ([#502](https://github.com/go-surreal/som/pull/502))
- Run against surrealdb v3.1.5 ([#505](https://github.com/go-surreal/som/pull/505))

### Fixed

- Circumvent surrealdb.go live query race condition ([#488](https://github.com/go-surreal/som/pull/488))
- Bump surrealdb.go and drop live query race workaround ([#495](https://github.com/go-surreal/som/pull/495))
- Remove unnecessary som version checks ([#496](https://github.com/go-surreal/som/pull/496))

## [v0.18.0] - 2026-03-09


### Added

- Add support for client-side transactions ([#463](https://github.com/go-surreal/som/pull/463))
- Surrealdb v3 structured server errors ([#465](https://github.com/go-surreal/som/pull/465))
- Raw queries ([#466](https://github.com/go-surreal/som/pull/466))
- Range query for string id nodes ([#471](https://github.com/go-surreal/som/pull/471))
- Where field value between ([#472](https://github.com/go-surreal/som/pull/472))
- Verify db version when connecting ([#474](https://github.com/go-surreal/som/pull/474))

### Changed

- Remove db client wrapper & cbor round-trip ([#468](https://github.com/go-surreal/som/pull/468))

### Fixed

- Surrealdb v3.0.2 query search same field  ([#462](https://github.com/go-surreal/som/pull/462))

### Documentation

- Prepare release v0.18.0 ([#475](https://github.com/go-surreal/som/pull/475))

## [v0.17.0] - 2026-03-01


### Added

- Add count index to all tables ([#445](https://github.com/go-surreal/som/pull/445))
- Fulltext search any (or) match ([#449](https://github.com/go-surreal/som/pull/449))
- Add index rebuild to repo methods ([#451](https://github.com/go-surreal/som/pull/451))
- Add insert method to repo for bulk create ([#452](https://github.com/go-surreal/som/pull/452))
- Add support for time.Month and time.Weekday ([#448](https://github.com/go-surreal/som/pull/448))
- Support proper struct sub-filter queries ([#453](https://github.com/go-surreal/som/pull/453))
- Add semver type with proper query filters ([#457](https://github.com/go-surreal/som/pull/457))
- Allow explicit custom field names via go tag ([#456](https://github.com/go-surreal/som/pull/456))
- Add geo types support ([#404](https://github.com/go-surreal/som/pull/404))
- Require surrealdb.go client version 1.3.0 ([#458](https://github.com/go-surreal/som/pull/458))

### Fixed

- Remove live query variable support workaround ([#446](https://github.com/go-surreal/som/pull/446))
- Align query rename offset to start ([#447](https://github.com/go-surreal/som/pull/447))
- Datetime cbor handling ([#450](https://github.com/go-surreal/som/pull/450))
- Properly align go tags syntax ([#455](https://github.com/go-surreal/som/pull/455))

### Documentation

- Update deprecated and add missing doc files ([#460](https://github.com/go-surreal/som/pull/460))

## [v0.16.0] - 2026-02-23


### Added

- Add support for complex id types ([#430](https://github.com/go-surreal/som/pull/430))
- Optimise cli command and flags ([#442](https://github.com/go-surreal/som/pull/442))
- Complex id range queries  ([#443](https://github.com/go-surreal/som/pull/443))

### Changed

- SubPkg to relativePkgPath ([#440](https://github.com/go-surreal/som/pull/440))
- Overhaul parser setup ([#441](https://github.com/go-surreal/som/pull/441))

## [v0.15.0] - 2026-02-21


### Added

- Optimise table id handling ([#427](https://github.com/go-surreal/som/pull/427))
- String record id with configurable generator ([#428](https://github.com/go-surreal/som/pull/428))
- Surrealdb v3.0 compatibility ([#438](https://github.com/go-surreal/som/pull/438))

### Changed

- Code cleanup & more tests ([#426](https://github.com/go-surreal/som/pull/426))
- Optimise go template handling ([#433](https://github.com/go-surreal/som/pull/433))

## [v0.14.0] - 2026-01-30


### Added

- Align query builder wording with official query lang ([#418](https://github.com/go-surreal/som/pull/418))
- Query fetch distinct field values ([#422](https://github.com/go-surreal/som/pull/422))
- Add lifecycle hooks to repo and model ([#423](https://github.com/go-surreal/som/pull/423))

### Fixed

- Expose errors properly ([#419](https://github.com/go-surreal/som/pull/419))
- Prevent nil pointer deref in repo read for missing records ([#421](https://github.com/go-surreal/som/pull/421))
- Slices with nil values ([#425](https://github.com/go-surreal/som/pull/425))

## [v0.13.0] - 2026-01-26


### Added

- Add most missing field filter functions ([#402](https://github.com/go-surreal/som/pull/402))
- Add support for github.com/gofrs/uuid ([#403](https://github.com/go-surreal/som/pull/403))
- Allow fetch for live queries ([#406](https://github.com/go-surreal/som/pull/406))
- Add opt-in context cache for database reads ([#409](https://github.com/go-surreal/som/pull/409))
- Optional soft-delete handling for nodes ([#410](https://github.com/go-surreal/som/pull/410))
- Generate wire setup if used by project ([#417](https://github.com/go-surreal/som/pull/417))

### Changed

- Cleanup query builder type ([#416](https://github.com/go-surreal/som/pull/416))

## [v0.12.0] - 2025-11-30


### Added

- Add indexing, unique and fulltext ([#153](https://github.com/go-surreal/som/pull/153))
- Add fulltext search features to query builder ([#397](https://github.com/go-surreal/som/pull/397))
- Add tempfiles statement to query builder ([#399](https://github.com/go-surreal/som/pull/399))

### Documentation

- Add fulltext search features ([#398](https://github.com/go-surreal/som/pull/398))

## [v0.11.0] - 2025-11-26


### Added

- Use readonly for created_at fields ([#385](https://github.com/go-surreal/som/pull/385))
- Make most field types sortable ([#387](https://github.com/go-surreal/som/pull/387))
- Add email field type ([#390](https://github.com/go-surreal/som/pull/390))
- Add iterator methods to query builder output ([#389](https://github.com/go-surreal/som/pull/389))
- Add smart password field type ([#300](https://github.com/go-surreal/som/pull/300))
- Optimistic locking feature for models ([#99](https://github.com/go-surreal/som/pull/99))

### Changed

- Remove unused fieldDef methods from fields ([#393](https://github.com/go-surreal/som/pull/393))

### Fixed

- Navigation and summary ([#391](https://github.com/go-surreal/som/pull/391))

### Documentation

- Initial documentation setup ([#305](https://github.com/go-surreal/som/pull/305))
- Add new features to gitbook ([#394](https://github.com/go-surreal/som/pull/394))

## [v0.10.0] - 2025-11-23


### Added

- Empty filter, schema and marshal fixup, testing ([#373](https://github.com/go-surreal/som/pull/373))

### Changed

- Better conv layer with direct cbor (un)marshal ([#374](https://github.com/go-surreal/som/pull/374))
- Cleanup legacy conv methods ([#381](https://github.com/go-surreal/som/pull/381))

### Fixed

- Schema generation for deeply nested fields ([#382](https://github.com/go-surreal/som/pull/382))

## [v0.9.0] - 2025-11-18


### Added

- Switch driver from sdbc to official surrealdb go client ([#371](https://github.com/go-surreal/som/pull/371))

### Fixed

- Main workflow wrong working dir ([#362](https://github.com/go-surreal/som/pull/362))

## [v0.8.0] - 2025-05-21


### Added

- Make som a compile-time only dependency ([#359](https://github.com/go-surreal/som/pull/359))
- Move cmd to root ([#361](https://github.com/go-surreal/som/pull/361))

### Fixed

- Go.mod commands fail with embedded template files ([#357](https://github.com/go-surreal/som/pull/357))
- Make config type alias to sdbc ([#358](https://github.com/go-surreal/som/pull/358))

## [v0.7.1] - 2025-05-04


### Added

- Db 2.0 enforced relations, enum string literals & more funcs ([#328](https://github.com/go-surreal/som/pull/328))

## [v0.7.0] - 2024-12-11


### Added

- Support for surrealdb 2.0 ([#321](https://github.com/go-surreal/som/pull/321))

## [v0.6.4] - 2024-11-30


### Fixed

- Testcontainers terminate does not ignore unnecessary error ([#323](https://github.com/go-surreal/som/pull/323))

## [v0.6.3] - 2024-09-05


### Fixed

- Codegen for (no-)pointer node/edge types ([#316](https://github.com/go-surreal/som/pull/316))

## [v0.6.2] - 2024-08-28


### Changed

- Hide tests in internal package ([#310](https://github.com/go-surreal/som/pull/310))
- Sdbc.ID to som.ID ([#311](https://github.com/go-surreal/som/pull/311))

## [v0.6.1] - 2024-08-27


### Fixed

- Generated files not written correctly ([#309](https://github.com/go-surreal/som/pull/309))

## [v0.6.0] - 2024-08-27


### Added

- Filter by direct field to field comparisons ([#302](https://github.com/go-surreal/som/pull/302))
- Add DescribeWithVars and Debug methods to query builder ([#304](https://github.com/go-surreal/som/pull/304))
- Add missing filter functions ([#303](https://github.com/go-surreal/som/pull/303))

### Changed

- Static param name in generated code ([#301](https://github.com/go-surreal/som/pull/301))

## [v0.5.0] - 2024-08-14


### Added

- Add table types to generated schema definition ([#264](https://github.com/go-surreal/som/pull/264))
- Optimise schema definition for timestamps ([#265](https://github.com/go-surreal/som/pull/265))
- Update sdbc & switch to cbor protocol ([#289](https://github.com/go-surreal/som/pull/289))
- Support time.duration type ([#164](https://github.com/go-surreal/som/pull/164))
- Implement filter functions and optimisation ([#291](https://github.com/go-surreal/som/pull/291))
- In-memory file system for codegen ([#293](https://github.com/go-surreal/som/pull/293))
- Check go, som and sdbc version before generate ([#266](https://github.com/go-surreal/som/pull/266))

### Changed

- Optimise cli setup ([#221](https://github.com/go-surreal/som/pull/221))
- Code cleanup ([#292](https://github.com/go-surreal/som/pull/292))
- Embedded gen files handling ([#298](https://github.com/go-surreal/som/pull/298))

### Fixed

- Hardcoded import path in embedded file ([#263](https://github.com/go-surreal/som/pull/263))
- Embedded file uses local import ([#296](https://github.com/go-surreal/som/pull/296))

## [v0.4.0] - 2024-05-02


### Added

- Implement live count method ([#234](https://github.com/go-surreal/som/pull/234))
- Generate comments for repo methods ([#247](https://github.com/go-surreal/som/pull/247))

### Fixed

- Filter for live queries ([#259](https://github.com/go-surreal/som/pull/259))

## [v0.3.0] - 2023-11-27


### Added

- Add support for byte and []byte fields ([#119](https://github.com/go-surreal/som/pull/119))
- Add support for net/url.URL type ([#163](https://github.com/go-surreal/som/pull/163))
- Add support for missing native numeric types ([#174](https://github.com/go-surreal/som/pull/174))
- Implement refresh method for repos ([#237](https://github.com/go-surreal/som/pull/237))

### Changed

- Reduce codegen by making query builder generic ([#245](https://github.com/go-surreal/som/pull/245))
- Move generic crud methods to embed ([#246](https://github.com/go-surreal/som/pull/246))

## [v0.2.0] - 2023-11-13


### Fixed

- Dependabot bump google.golang.org/grpc to 1.57.1 ([#231](https://github.com/go-surreal/som/pull/231))
- Dependabot bump github.com/docker/docker to 24.0.7 ([#233](https://github.com/go-surreal/som/pull/233))
- Nil deref due to invalid use of anonymous struct field type ([#239](https://github.com/go-surreal/som/pull/239))
- Edge field conversion ([#241](https://github.com/go-surreal/som/pull/241))
- Remove duplicates from generated code ([#240](https://github.com/go-surreal/som/pull/240))
- Generated names for types & fields ([#238](https://github.com/go-surreal/som/pull/238))
- Generated filters missing methods for edges ([#243](https://github.com/go-surreal/som/pull/243))

## [v0.1.2] - 2023-10-16


### Fixed

- Timestamps zero value deref nil pointer ([#227](https://github.com/go-surreal/som/pull/227))

## [v0.1.1] - 2023-10-09


### Added

- Add client options and pass to sdbc ([#222](https://github.com/go-surreal/som/pull/222))

## [v0.1.0] - 2023-09-17


### Added

- Create nodes with ulid as primary key ([#204](https://github.com/go-surreal/som/pull/204))

### Changed

- Move repo to go-surreal org ([#194](https://github.com/go-surreal/som/pull/194))
- Move sdbc code to own repo ([#203](https://github.com/go-surreal/som/pull/203))
- Move timestamp handling to database schema layer ([#197](https://github.com/go-surreal/som/pull/197))

### Fixed

- Codegen bugs & streamline tests ([#210](https://github.com/go-surreal/som/pull/210))

### Documentation

- Update readme ([#201](https://github.com/go-surreal/som/pull/201))
- Update readme & add ideas ([#160](https://github.com/go-surreal/som/pull/160))

## [v0.0.11] - 2023-09-05


### Added

- Add surrealdb strict-mode compatibility ([#175](https://github.com/go-surreal/som/pull/175))
- Implement custom surrealdb client ([#152](https://github.com/go-surreal/som/pull/152))
- Optimise sdbc close & cleanup handling ([#191](https://github.com/go-surreal/som/pull/191))
- Add support for surrealdb beta.11 ([#200](https://github.com/go-surreal/som/pull/200))
- Add support for live queries ([#167](https://github.com/go-surreal/som/pull/167))

### Changed

- Let conv functions use pointers ([#198](https://github.com/go-surreal/som/pull/198))

### Fixed

- Apply schema with transaction is broken ([#186](https://github.com/go-surreal/som/pull/186))
- Apply latest surrealdb nightly changes ([#187](https://github.com/go-surreal/som/pull/187))
- Enum & struct fields handling ([#199](https://github.com/go-surreal/som/pull/199))

## [v0.0.10] - 2023-08-04


### Added

- Add transaction wrapping to schema definition ([#156](https://github.com/go-surreal/som/pull/156))
- Add assert for uuid fields to schema ([#162](https://github.com/go-surreal/som/pull/162))
- Add id field definition to generated schema ([#165](https://github.com/go-surreal/som/pull/165))
- Add async methods to query builder ([#171](https://github.com/go-surreal/som/pull/171))

### Changed

- Simplify conv of uuid type ([#166](https://github.com/go-surreal/som/pull/166))
- Move static code to embed files ([#170](https://github.com/go-surreal/som/pull/170))
- Optimise query builder string concat ([#173](https://github.com/go-surreal/som/pull/173))

### Fixed

- Count query and result mapping ([#130](https://github.com/go-surreal/som/pull/130))
- Unexpected db response structure for unmarshal ([#131](https://github.com/go-surreal/som/pull/131))
- Dependabot alert by bumping google.golang.org/grpc to 1.53.0 ([#145](https://github.com/go-surreal/som/pull/145))
- Tests ([#157](https://github.com/go-surreal/som/pull/157))
- Schema needs to set option<?> types for nilable fields ([#161](https://github.com/go-surreal/som/pull/161))

## [v0.0.9] - 2023-04-09


### Added

- Utilize new features of surrealdb.go client update v0.2.0 ([#125](https://github.com/go-surreal/som/pull/125))
- Update for surrealdb version beta.9 ([#129](https://github.com/go-surreal/som/pull/129))

### Fixed

- Dependabot alert by updating opencontainers/runc to 1.1.5 ([#123](https://github.com/go-surreal/som/pull/123))

## [v0.0.8] - 2023-03-29


### Added

- Make generated code mockable by using interfaces ([#112](https://github.com/go-surreal/som/pull/112))
- Small gen updates, more tests and examples ([#106](https://github.com/go-surreal/som/pull/106))
- Move lib package into generated sources via embed ([#118](https://github.com/go-surreal/som/pull/118))

### Documentation

- Update readme & faq ([#124](https://github.com/go-surreal/som/pull/124))

## [v0.0.7] - 2023-02-13


### Fixed

- Record links ([#98](https://github.com/go-surreal/som/pull/98))

### Documentation

- Add ideas document ([#92](https://github.com/go-surreal/som/pull/92))

## [v0.0.6] - 2023-01-24


### Added

- Implement sub-queries for nodes and edges ([#84](https://github.com/go-surreal/som/pull/84))

## [v0.0.5] - 2023-01-19


### Added

- Add asserts to database schema definition ([#82](https://github.com/go-surreal/som/pull/82))
- Schema assert for enum values ([#83](https://github.com/go-surreal/som/pull/83))

### Documentation

- Update readme ([#72](https://github.com/go-surreal/som/pull/72))

## [v0.0.4] - 2022-12-29


### Added

- Implement query describe method ([#74](https://github.com/go-surreal/som/pull/74))
- Make query builder reusable ([#76](https://github.com/go-surreal/som/pull/76))
- Implement auto timestamps on create/update ([#79](https://github.com/go-surreal/som/pull/79))

### Changed

- Provide separate create-with-id method ([#80](https://github.com/go-surreal/som/pull/80))

### Fixed

- Query variable index breaks after 26 variables ([#75](https://github.com/go-surreal/som/pull/75))

## [v0.0.3] - 2022-12-26


### Added

- Implement select, update and delete ([#57](https://github.com/go-surreal/som/pull/57))
- Add "do not edit" comment to generated code ([#63](https://github.com/go-surreal/som/pull/63))
- Allow custom record ids ([#68](https://github.com/go-surreal/som/pull/68))
- Support pointer fields ([#64](https://github.com/go-surreal/som/pull/64))
- Generate database schema for strict usage ([#44](https://github.com/go-surreal/som/pull/44))

### Changed

- Update database interface ([#58](https://github.com/go-surreal/som/pull/58))
- Code generator for better extensibility ([#66](https://github.com/go-surreal/som/pull/66))
- Code generation streamlining ([#67](https://github.com/go-surreal/som/pull/67))

## [v0.0.2] - 2022-12-08


### Added

- Enforce conventional commits spec via PR workflow ([#37](https://github.com/go-surreal/som/pull/37))

### Fixed

- Golangci-lint issues ([#52](https://github.com/go-surreal/som/pull/52))
- Lint example/gen and handle existing issues ([#56](https://github.com/go-surreal/som/pull/56))

### Documentation

- Update README ([#55](https://github.com/go-surreal/som/pull/55))

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

