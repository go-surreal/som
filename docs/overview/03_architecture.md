# Architecture

This page describes how SOM works internally.

## Generation Pipeline

SOM uses a code generation approach to provide type-safe database access. The generation workflow:

```
Input Models (Go structs)
        │
        ▼
┌───────────────────┐
│      Parser       │  Analyzes Go source using gotype
│   (core/parser)   │  Identifies: Nodes, Edges, Structs, Enums
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│  Code Generator   │  Uses jennifer for Go code generation
│  (core/codegen)   │  Produces type-safe builders
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│   Output Files    │  Complete database access layer
│    (gen/som/)     │  Ready to use in your application
└───────────────────┘
```

### 1. Parsing Phase

The parser analyzes your Go source files to identify:

- **Nodes**: Structs embedding `som.Node`
- **Edges**: Structs embedding `som.Edge`
- **Structs**: Regular structs used as fields (for nested types)
- **Enums**: Types implementing `som.Enum` interface
- **Fields**: All fields with their types, tags, and constraints

### 2. Code Generation Phase

Using the [jennifer](https://github.com/dave/jennifer) library, SOM generates idiomatic Go code:

- Static library code (filters, builders, utilities)
- Per-model repositories and query builders
- Type converters for database serialization
- Filter and sort definitions for each field type

### 3. Output Phase

Generated files are written to the output directory with their own `go.mod`, creating a self-contained package.

## Generated Code Structure

```
gen/som/
├── som.base.go           # ID, Node, Edge, Timestamps types
├── som.client.go         # Client implementation
├── som.interfaces.go     # Repository interfaces
├── som.node.go           # Base repository pattern
├── som.schema.go         # Schema utilities
│
├── repo/                 # Repository implementations
│   └── node.{model}.go   # Per-model: Create, Read, Update, Delete, Query
│
├── query/                # Query builders
│   ├── builder.go        # Generic builder with all methods
│   ├── query.go          # Query execution, async, live
│   └── node.{model}.go   # Per-model query factory
│
├── where/                # Filter definitions
│   ├── where.go          # Initialization
│   ├── node.{model}.go   # Per-model field filters
│   └── object.{struct}.go # Nested struct filters
│
├── by/                   # Sort definitions
│   └── node.{model}.go   # Per-model sortable fields
│
├── with/                 # Fetch/eager loading
│   └── node.{model}.go   # Per-model fetchable relations
│
├── conv/                 # Type converters
│   ├── node.{model}.go   # Model to/from database format
│   └── edge.{edge}.go    # Edge conversions
│
├── relate/               # Edge creation
│   └── edge.{edge}.go    # Per-edge RELATE builders
│
├── internal/
│   ├── lib/              # Filter and sort implementation
│   ├── cbor/             # CBOR marshaling utilities
│   └── types/            # DateTime, Duration, UUID wrappers
│
└── constant/             # Math constants (PI, E, etc.)
```

## Key Components

### Nodes

Nodes represent database records. Any struct embedding `som.Node` becomes a SurrealDB table:

```go
type User struct {
    som.Node        // Provides ID field
    som.Timestamps  // Optional: CreatedAt, UpdatedAt
    Name string
}
// Creates table: user
```

### Edges

Edges represent graph relationships between nodes:

```go
type Follows struct {
    som.Edge        // Provides ID, In, Out
    Since time.Time // Edge metadata
}
// Creates edge table: follows
// RELATE user:a->follows->user:b
```

### Repositories

Generated for each model with full CRUD operations:

```go
type UserRepo interface {
    Create(ctx, *User) error
    CreateWithID(ctx, id string, *User) error
    Read(ctx, *som.ID) (*User, bool, error)
    Update(ctx, *User) error
    Delete(ctx, *User) error
    Refresh(ctx, *User) error
    Query() Builder[User, ConvUser]
}
```

### Query Builder

Fluent API with method chaining:

```go
Builder[M, C]
├── Filter(filters...)     // WHERE conditions
├── Order(sorts...)        // ORDER BY
├── OrderRandom()          // ORDER RAND()
├── Offset(n)              // OFFSET
├── Limit(n)               // LIMIT
├── Fetch(relations...)    // FETCH (eager load)
├── Timeout(duration)      // Execution timeout
├── Parallel(bool)         // Parallel execution
│
├── All(ctx)               // Get all results
├── First(ctx)             // Get first result
├── One(ctx)               // Get exactly one
├── Count(ctx)             // Count results
├── Exists(ctx)            // Check existence
├── Live(ctx)              // Stream changes
│
└── *Async variants        // AllAsync, FirstAsync, etc.
```

### Type Converters

Handle transformation between Go types and SurrealDB format:

```go
// Generated in conv/
func ToUser(model *model.User) *convUser { ... }
func FromUser(conv *convUser) *model.User { ... }
```

## Database Communication

### CBOR Protocol

SOM uses CBOR (Concise Binary Object Representation) for SurrealDB communication. Custom types use CBOR tags:

| Type | CBOR Tag | Format |
|------|----------|--------|
| DateTime | 12 | `[unix_seconds, nanoseconds]` |
| Duration | 14 | Nanoseconds as int64 |
| UUID | 37 | 16-byte binary |

## ID Handling

SOM provides utilities for working with SurrealDB record IDs:

```go
// Create a record ID
id := som.NewRecordID("user", "abc123")

// Create a pointer to record ID
idPtr := som.MakeID("user", "abc123")

// Reference a table (for auto-generated IDs)
table := som.Table("user")
```

Record IDs are automatically generated as ULIDs when using `Create()`.
