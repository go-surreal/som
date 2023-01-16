# ROADMAP

## Before v0.1.0 (first "somewhat stable" non-pre-release)

- [ ] Implement sub-queries for node and edge types.
- [ ] Add `som.SoftDelete` type with `DeletedAt` timestamp and automated handling throughout som.
- [ ] Mark fetched sub-nodes as "invalid to be saved"? (#25)
- [ ] Consider reserved (query) keywords. (#18)
- [ ] Check for possible security vulnerabilities.
- [ ] Choose proper licensing for the project. (#11)

## After v0.1.0

- [ ] Provide `WithInfo` method.
- [ ] Add support for `[]byte` (and `byte`?) type.
- [ ] How to handle data migrations? (#22)
- [ ] Setup golangci-lint with proper config. (#7)
- [ ] Support (deeply) nested slices? (needed?)
- [ ] Cleanup naming conventions. (#24)
- [ ] Code comments and documentation. (#9)
- [ ] Write tests. (#8)
- [ ] Generate `sommock` package for easy mocking of the underlying database client.
- [ ] Make casing of database field names configurable.
- [ ] Switch the source code parser to support generics.
- [ ] Add `som.Edge[I, O any]` for defining edges more clearly and without tags (requires generics parser).
- [ ] Support transactions.
- [ ] Distinct results (https://stackoverflow.com/questions/74326176/surrealdb-equivalent-of-select-distinct).
- [ ] Integrate external APIs (GraphQL) into the db access layer?
- [ ] Support (deeply) nested slices? (needed?)
- [ ] Unique relations (`DEFINE INDEX __som_unique_relation ON TABLE member_of COLUMNS in, out UNIQUE;`)

## Nice to have (v0.x.x)?

- [ ] Add new data type "password" with automatic handling of encryption with salt. (#16)
- [ ] Add data type "email" as alias for string that adds database assertion.
    - Or provide an API to add custom assertions for types (especially string).
- [ ] Add performance benchmarks (and possible optimizations due to it).

```sql
DEFINE TABLE user SCHEMAFULL 
        PERMISSIONS NONE;
DEFINE FIELD username ON TABLE user
        TYPE string
        ASSERT string::length($value) >= 4
        ASSERT string::length($value) <= 8;
DEFINE FIELD password ON TABLE user
        PERMISSIONS
                FOR SELECT NONE
        TYPE string;
DEFINE FIELD email ON TABLE user
        TYPE string
        ASSERT is::email($value);
DEFINE FIELD num ON TABLE user
        VALUE 42;
```
