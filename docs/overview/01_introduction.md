# Introduction

SOM (SurrealDB Object Mapper) is an ORM and query builder for [SurrealDB](https://surrealdb.com/) with built-in model mapping and type-safe query operation generator. It provides an easy and sophisticated database access layer for Go applications.

## What is SurrealDB?

SurrealDB is a cutting-edge database system that offers a SQL-style query language with real-time queries and efficient related data retrieval. It supports both schema-full and schema-less data handling.

With its full graph database functionality, SurrealDB enables advanced querying and analysis by allowing records (or vertices) to be connected with edges, each with its own properties and metadata. This facilitates multi-table, multi-depth document retrieval without complex JOINs, all within the database.

*(Information extracted from the [official homepage](https://surrealdb.com))*.

## Why SOM?

Working directly with SurrealDB in Go requires manual query building, type conversion, and result mapping. SOM eliminates this boilerplate by generating type-safe code from your Go struct definitions.

**SOM provides:**

- **Type-safe queries**: Compile-time checked queries prevent runtime errors. Field names, types, and operations are all validated by the Go compiler.
- **Code generation**: Automatically generates repository and query builder code from your Go models. No reflection at runtime.
- **Native Go experience**: Work with your domain models directly. SOM handles conversion between Go types and SurrealDB's CBOR protocol.
- **Graph support**: First-class support for SurrealDB's graph capabilities via typed edges with metadata.
- **Real-time queries**: Built-in support for SurrealDB's live query feature with type-safe event handling.
- **Async operations**: All query methods have async variants for concurrent database access.

## How It Works

1. **Define models** as Go structs embedding `som.Node` (for records) or `som.Edge` (for relationships)
2. **Run the generator** to produce type-safe database access code
3. **Use the generated client** with full IDE autocompletion and compile-time safety

```go
// 1. Define your model
type User struct {
    som.Node
    Name  string
    Email string
}

// 2. Generate code: som gen ./model ./gen/som

// 3. Use the generated client
users, _ := client.UserRepo().Query().
    Filter(where.User.Email.Contains("@example.com")).
    All(ctx)
```

## Project Status

SOM is currently in **early development** and considered experimental. While core functionality works, the API may change between versions. See the [FAQ](../appendix/01_faq.md) for details on production readiness.

## Dependencies

SOM uses [sdbc](https://github.com/go-surreal/sdbc), a custom SurrealDB client optimized for code generation patterns. This provides CBOR protocol support and features tailored for SOM's needs. The official SurrealDB Go client is not currently supported.
