# Architecture

This page describes how SOM works internally.

## Generation Pipeline

SOM uses a code generation approach to provide type-safe database access. The generation workflow consists of these steps:

1. **Parse**: Analyze Go source files to identify models (structs embedding `som.Node` or `som.Edge`)
2. **Generate Static Code**: Output library code for filters, builders, and utilities
3. **Generate Model Code**: Create model-specific queries, repositories, and converters
4. **Write Output**: Save all generated files to the specified output directory

## Key Components

### Nodes

Nodes represent database records. Any struct that embeds `som.Node` is treated as a database table.

### Edges

Edges represent relationships between nodes. Structs embedding `som.Edge` create graph edges with `in` and `out` fields pointing to connected nodes.

### Repositories

Generated repositories provide CRUD operations for each model:
- `Create` - Insert new records
- `Read` - Fetch by ID
- `Update` - Modify existing records
- `Delete` - Remove records
- `Query` - Build complex queries

### Query Builder

The generated query builder provides a fluent API for constructing type-safe queries with filters, ordering, and pagination.

## Database Client

SOM uses a custom SurrealDB client called [sdbc](https://github.com/go-surreal/sdbc) rather than the official client. This provides optimized CBOR serialization and features tailored for code generation.
