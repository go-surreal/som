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
- [ ] add support for maps as node fields

## References

- https://surrealdb.com/docs
- https://entgo.io
- https://github.com/d-tsuji/awesome-go-orms
- https://github.com/doug-martin/goqu
- https://github.com/sharovik/orm
