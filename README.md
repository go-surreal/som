# SDB

A type-safe database access layer with code generation for [SurrealDB](https://surrealdb.com).

## Generate the client code

```
go run github.com/marcbinz/sdb/cmd/sdbgen@latest <input_dir> <output_dir>
```

## Roadmap

### TODOs:

- [ ] generate `sdbmock` package for easy mocking of the underlying database client
- [ ] add performance benchmarks (and possible optimizations due to it)
- [ ] make casing of database fields configurable?
- [ ] add support for maps as node fields?
- [ ] add support for `[]byte` (and `byte`?) type
- [ ] Replace `strcase.ToLowerCamel` with `strcase.ToLowerCamel`
- [ ] Switch the source code parsing to support generics
- [ ] add `som.Edge[I, O any]` for defining edges more clearly and without tags (requires generics parser)
- [ ] Make query builder not use pointers, so partial builds and usages are working
- [ ] Fix query variable index `rune` not suited for > 26 due to invalid char as variable key?

## FAQ

### Why are maps not supported?

- With the schemaless database this would be possible.
- Currently, the focus is on structured and deterministic data.
- Might be added in the future though.

## References

- https://surrealdb.com/docs
- https://entgo.io
- https://github.com/d-tsuji/awesome-go-orms
- https://github.com/doug-martin/goqu
- https://github.com/sharovik/orm
