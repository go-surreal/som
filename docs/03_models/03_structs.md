# Embedded Structs

SOM supports embedding regular Go structs within your models for better code organization and reuse.

## Basic Embedding

Embed any struct within a Node as a field:

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
    Address Address  // Embedded struct - stored as nested object
}
```

In SurrealDB, this creates a nested object:
```json
{
  "id": "user:abc",
  "name": "Alice",
  "address": {
    "street": "123 Main St",
    "city": "Berlin",
    "country": "Germany",
    "zip_code": "10115"
  }
}
```

## Pointer Embedding (Optional)

Use pointers for optional embedded structs:

```go
type User struct {
    som.Node

    Name    string
    Address *Address  // Optional, can be nil
}
```

This allows:
```go
// User without address
user := &model.User{Name: "Bob"}

// User with address
user := &model.User{
    Name: "Alice",
    Address: &model.Address{
        City: "Berlin",
    },
}
```

## Filtering on Nested Fields

SOM generates filters for nested struct fields:

```go
// Filter by nested field
users, err := client.UserRepo().Query().
    Filter(where.User.Address.City.Equal("Berlin")).
    All(ctx)

// Multiple nested filters
users, err := client.UserRepo().Query().
    Filter(
        where.User.Address.Country.Equal("Germany"),
        where.User.Address.City.Contains("Ber"),
    ).
    All(ctx)
```

## Deeply Nested Structs

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

Filter on deeply nested fields:
```go
// Filter by latitude
users, err := client.UserRepo().Query().
    Filter(where.User.Address.Coordinates.Lat.GreaterThan(52.0)).
    All(ctx)
```

## Multiple Embeddings

Embed multiple structs for composition:

```go
type Metadata struct {
    Version   int
    Source    string
    CreatedBy string
}

type Settings struct {
    Theme       string
    Language    string
    Timezone    string
}

type User struct {
    som.Node
    som.Timestamps

    Name     string
    Address  Address
    Metadata Metadata
    Settings *Settings  // Optional
}
```

## Struct Reuse

The same struct can be used in multiple nodes:

```go
type ContactInfo struct {
    Email string
    Phone string
}

type User struct {
    som.Node
    Name    string
    Contact ContactInfo
}

type Company struct {
    som.Node
    Name    string
    Contact ContactInfo  // Same struct
}
```

SOM generates separate filter definitions for each usage:
```go
where.User.Contact.Email.Equal("...")
where.Company.Contact.Email.Equal("...")
```

## Struct with Special Types

Embedded structs can contain any supported type:

```go
type Profile struct {
    Bio       string
    Website   *url.URL
    Birthday  *time.Time
    AvatarID  uuid.UUID
    Tags      []string
    Verified  bool
}

type User struct {
    som.Node
    Profile Profile
}
```

## Anonymous Embedding

Go's anonymous embedding promotes fields to the parent:

```go
type Auditable struct {
    CreatedBy string
    UpdatedBy string
}

type Document struct {
    som.Node
    Auditable  // Anonymous - fields promoted

    Title   string
    Content string
}

// Access promoted fields
doc.CreatedBy = "alice"
```

## Limitations

- Structs cannot embed `som.Node` or `som.Edge` (those create separate tables)
- Circular struct references are not supported
- Very deep nesting may impact query performance
