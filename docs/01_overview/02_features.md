# Features

## Core Features

- **Code Generation**: Generate type-safe database access code from Go struct models
- **Query Builder**: Fluent API for building complex queries with compile-time type checking
- **Model Mapping**: Automatic conversion between Go structs and SurrealDB records
- **Relationship Support**: Native support for graph edges and record links

## Supported Data Types

### Primitive Types

- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint8`, `uint16`, `uint32`
- `float32`, `float64`
- `bool`
- `rune`
- `byte`, `[]byte`

### Time Types

- `time.Time`
- `time.Duration`

### Special Types

- `url.URL`
- `github.com/google/uuid.UUID`

### Custom Types

- `som.Enum` - For enumerated values

### Collections

- Slice types of all supported types
- Pointer versions of all types

## Current Limitations

### Unsupported Go Types

- `uint`, `uint64`, `uintptr` - Due to SurrealDB limitations with very large integers
- `complex64`, `complex128`
- `map` types (planned for future)
