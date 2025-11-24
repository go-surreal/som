# CLI Usage

The SOM CLI generates type-safe database access code from your Go models.

## Basic Command

```bash
go run github.com/go-surreal/som/cmd/som@latest gen <input_dir> <output_dir>
```

### Arguments

- `<input_dir>` - Directory containing your model files (structs with `som.Node` or `som.Edge`)
- `<output_dir>` - Directory where generated code will be written

### Example

```bash
# Generate from ./model to ./gen/som
go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som
```

## Installing the Binary

For faster execution, install the binary:

```bash
go install github.com/go-surreal/som/cmd/som@latest
```

Then run directly:

```bash
som gen ./model ./gen/som
```

## Output Structure

The generator creates the following structure:

```
gen/som/
├── som.base.go         # Base types (Node, Edge, ID utilities)
├── by/                 # Order/sort helpers
│   └── <model>.go      # Per-model ordering
├── constant/           # Generated constants
├── conv/               # Model converters (internal ↔ external)
│   └── <model>.go      # Per-model converter
├── internal/           # Internal utilities
│   ├── cbor/           # CBOR encoding helpers
│   ├── lib/            # Filter/query library
│   └── types/          # Custom types (DateTime, Duration, UUID)
├── query/              # Query builder implementations
│   └── <model>.go      # Per-model query builder
├── relate/             # Edge relationship builders
│   └── <edge>.go       # Per-edge relate builder
├── repo/               # Repository implementations
│   └── <model>.go      # Per-model repository
├── where/              # Filter condition builders
│   └── <model>.go      # Per-model filters
└── with/               # Fetch/eager loading helpers
    └── <model>.go      # Per-model fetch paths
```

## Generated Packages

### `som` (root)

Core types and client:

- `Node` - Base type for database records
- `Edge` - Base type for relationships
- `ID` - Record ID type alias
- `NewRecordID()`, `MakeID()`, `Table()` - ID utilities
- `Timestamps` - Auto-managed timestamp fields
- `Enum`, `Email`, `Password`, `SemVer` - Special types

### `repo`

Repository implementations with CRUD operations:

```go
type UserRepo interface {
    Create(ctx, user) error
    CreateWithID(ctx, id, user) error
    Read(ctx, id) (*model.User, bool, error)
    Update(ctx, user) error
    Delete(ctx, user) error
    Refresh(ctx, user) error
    Query() query.Builder[model.User, conv.User]
}
```

### `query`

Fluent query builders:

```go
client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    Order(by.User.Name.Asc()).
    Limit(10).
    All(ctx)
```

### `where`

Type-safe filter conditions:

```go
where.User.Email.Contains("@example.com")
where.User.Age.GreaterThan(18)
where.User.CreatedAt.After(lastWeek)
```

### `by`

Ordering definitions:

```go
by.User.Name.Asc()
by.User.CreatedAt.Desc()
```

### `with`

Fetch paths for eager loading:

```go
query.Fetch(with.User.Posts...)
query.Fetch(with.Post.Author...)
```

### `conv`

Internal converters between model and database representations. Not typically used directly.

### `relate`

Edge relationship builders:

```go
client.FollowsRepo().Relate().
    From(alice).
    To(bob).
    Create(ctx, follows)
```

## Version Pinning

Pin to a specific version for reproducible builds:

```bash
go run github.com/go-surreal/som/cmd/som@v0.10.0 gen ./model ./gen/som
```

## Regenerating Code

Run the generator whenever you:

- Add or remove model fields
- Create new Node or Edge types
- Rename models or fields
- Change field types
- Add or remove struct tags

## Docker Usage

Run from Docker for consistent builds:

```bash
docker run --rm -v $(pwd):/app -w /app golang:1.23 \
    go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som
```

## Troubleshooting

### Models not detected

Ensure structs embed `som.Node` or `som.Edge`:

```go
type User struct {
    som.Node  // Required!
    Name string
}
```

### Import errors

After model changes, regenerate and run:

```bash
som gen ./model ./gen/som
go mod tidy
```

### Conflicting types

If you get type conflicts, ensure your model package has a unique name and doesn't shadow generated packages.
