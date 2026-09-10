# Primitive Types

SOM supports most Go primitive types for use in your models.

## Strings

```go
type User struct {
    som.Node[som.ULID]

    Name        string
    Description string
}
```

## Integers

Supported signed integers:

```go
type Metrics struct {
    som.Node[som.ULID]

    Count8  int8
    Count16 int16
    Count32 int32
    Count64 int64
    Count   int
}
```

Supported unsigned integers:

```go
type Metrics struct {
    som.Node[som.ULID]

    Size8  uint8
    Size16 uint16
    Size32 uint32
}
```

> **Note**: `uint`, `uint64`, and `uintptr` are not currently supported due to SurrealDB limitations with very large integers.

## Floating Point

```go
type Measurement struct {
    som.Node[som.ULID]

    Value32 float32
    Value64 float64
}
```

## Booleans

```go
type User struct {
    som.Node[som.ULID]

    IsActive  bool
    IsAdmin   bool
}
```

## Bytes

```go
type Document struct {
    som.Node[som.ULID]

    SingleByte byte
    Data       []byte
}
```

## Runes

```go
type Character struct {
    som.Node[som.ULID]

    Symbol rune
}
```

## Pointers

All primitive types support pointer versions for optional values:

```go
type User struct {
    som.Node[som.ULID]

    Name     string   // Required
    Nickname *string  // Optional, can be nil
    Age      *int     // Optional
}
```
