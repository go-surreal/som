# Embedded Structs

SOM supports embedding regular Go structs within your models for better code organization and reuse.

## Basic Embedding

You can embed any struct within a Node:

```go
type Address struct {
    Street  string
    City    string
    Country string
    ZipCode string
}

type User struct {
    som.Node

    Name    string
    Address Address  // Embedded struct
}
```

## Pointer Embedding

Use pointers for optional embedded structs:

```go
type User struct {
    som.Node

    Name    string
    Address *Address  // Optional, can be nil
}
```

## Multiple Embeddings

Embed multiple structs for composition:

```go
type Metadata struct {
    Version   int
    Source    string
}

type User struct {
    som.Node
    som.Timestamps

    Name     string
    Address  Address
    Metadata Metadata
}
```

## Nested Structs

Structs can be nested multiple levels deep:

```go
type Coordinates struct {
    Lat float64
    Lng float64
}

type Address struct {
    Street      string
    City        string
    Coordinates Coordinates
}

type User struct {
    som.Node
    Address Address
}
```

The generated query builder supports filtering on nested fields.
