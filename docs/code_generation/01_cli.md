# CLI Usage

The SOM CLI generates type-safe database access code from your Go models.

## Basic Command

```bash
go run github.com/go-surreal/som@latest -i <input_dir>
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-i`, `--in` | Input model directory (relative to go.mod) | *(none, omit for init mode)* |
| `-o`, `--out` | Output directory | `gen/som` |
| `-v`, `--verbose` | Show detailed generation progress | `false` |
| `--dry` | Simulate generation without writing files | `false` |
| `--no-check` | Skip version compatibility checks | `false` |
| `--no-count-index` | Disable automatic COUNT index generation | `false` |
| `--wire` | Override wire generation: `"no"`, `"google"`, or `"goforj"` | auto-detect |

### Example

```bash
# Generate from ./model to ./gen/som (default output)
go run github.com/go-surreal/som@latest -i ./model

# Generate to custom output directory
go run github.com/go-surreal/som@latest -i ./model -o ./db/generated
```

## Initialization Mode

When `--in` is omitted, SOM generates only static base files. This is useful for bootstrapping a new project so types like `som.Node` are available before you define models:

```bash
go run github.com/go-surreal/som@latest
```

## Installing the Binary

For faster execution, install the binary:

```bash
go install github.com/go-surreal/som@latest
```

Then run directly:

```bash
som -i ./model
```

## Output Structure

The generator creates the following structure:

```
gen/som/
├── som.base.go         # Base types (Node, Edge, ID utilities)
├── by/                 # Sort helpers
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
├── filter/             # Filter condition builders
│   └── <model>.go      # Per-model filters
└── with/               # Fetch/eager loading helpers
    └── <model>.go      # Per-model fetch paths
```

## Generated Packages

### `som` (root)

Core types and client:

- `Node[T]` - Base type for database records (generic over ID type)
- `Edge` - Base type for relationships
- `Timestamps` - Auto-managed timestamp fields
- `OptimisticLock` - Version-based conflict detection
- `SoftDelete` - Non-destructive deletion
- `Email`, `Password[A]`, `SemVer` - Special types
- `ULID`, `UUID`, `Rand` - ID types
- `ArrayID`, `ObjectID` - Complex ID markers

### `repo`

Repository implementations with CRUD operations:

```go
type UserRepo interface {
    Create(ctx, user) error
    CreateWithID(ctx, id, user) error
    Insert(ctx, users) error
    Read(ctx, id) (*model.User, bool, error)
    Update(ctx, user) error
    Delete(ctx, user) error
    Refresh(ctx, user) error
    Index() *index.User
    Query() query.Builder[model.User]
}
```

### `query`

Fluent query builders:

```go
client.UserRepo().Query().
    Where(filter.User.IsActive.IsTrue()).
    Order(by.User.Name.Asc()).
    Limit(10).
    All(ctx)
```

### `filter`

Type-safe filter conditions:

```go
filter.User.Email.Contains("@example.com")
filter.User.Age.GreaterThan(18)
filter.User.CreatedAt.After(lastWeek)
```

### `by`

Sort definitions:

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
go run github.com/go-surreal/som@v0.16.0 -i ./model
```

## Regenerating Code

Run the generator whenever you:

- Add or remove model fields
- Create new Node or Edge types
- Rename models or fields
- Change field types
- Add or remove struct tags

## Troubleshooting

### Models not detected

Ensure structs embed `som.Node[T]` or `som.Edge`:

```go
type User struct {
    som.Node[som.ULID]  // Required!
    Name string
}
```

### Import errors

After model changes, regenerate and run:

```bash
go run github.com/go-surreal/som@latest -i ./model
go mod tidy
```

### Conflicting types

If you get type conflicts, ensure your model package has a unique name and doesn't shadow generated packages.
