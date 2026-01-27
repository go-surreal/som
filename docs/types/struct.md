# Struct Type

The struct type handles embedded Go structs as nested objects within your models.

## Overview

| Property | Value |
|----------|-------|
| Go Type | Embedded struct / `*EmbeddedStruct` |
| Database Schema | `object` / `option<object>` |
| CBOR Encoding | Direct |
| Sortable | No (use nested field sorting) |

## Definition

Define structs in your model package and embed them:

```go
package model

// Embedded struct (not a Node, no ID)
type Address struct {
    Street  string
    City    string
    Country string
    ZipCode string
}

type Dimensions struct {
    Width  float64
    Height float64
    Depth  float64
}

type User struct {
    som.Node

    Name    string
    Address Address    // Required embedded struct
    Billing *Address   // Optional embedded struct
}

type Product struct {
    som.Node

    Name       string
    Dimensions Dimensions
}
```

> **Note**: Embedded structs do NOT embed `som.Node` - they are plain Go structs.

## Schema

Generated SurrealDB schema:

```surql
DEFINE FIELD address ON user TYPE object;
DEFINE FIELD address.street ON user TYPE string;
DEFINE FIELD address.city ON user TYPE string;
DEFINE FIELD address.country ON user TYPE string;
DEFINE FIELD address.zip_code ON user TYPE string;

DEFINE FIELD billing ON user TYPE option<object>;
DEFINE FIELD billing.street ON user TYPE option<string>;
DEFINE FIELD billing.city ON user TYPE option<string>;
DEFINE FIELD billing.country ON user TYPE option<string>;
DEFINE FIELD billing.zip_code ON user TYPE option<string>;
```

## Creating Struct Values

```go
user := &model.User{
    Name: "Alice",
    Address: model.Address{
        Street:  "123 Main St",
        City:    "New York",
        Country: "USA",
        ZipCode: "10001",
    },
}

// Optional struct
billing := &model.Address{
    Street:  "456 Business Ave",
    City:    "Chicago",
    Country: "USA",
    ZipCode: "60601",
}
user.Billing = billing
```

## Filter Operations

### Nested Field Filtering

Access nested fields using dot notation in the where clause:

```go
// Filter by nested field
filter.User.Address.City.Equal("New York")

// String operations on nested fields
filter.User.Address.Country.Equal("USA")
filter.User.Address.ZipCode.StartsWith("100")

// Numeric nested fields
filter.Product.Dimensions.Width.GreaterThan(10.0)
filter.Product.Dimensions.Height.LessThan(50.0)
```

### Complex Nested Filters

```go
// Multiple conditions on nested struct
query.Where(
    filter.User.Address.City.Equal("New York"),
    filter.User.Address.Country.Equal("USA"),
)
```

### Optional Struct Nil Checks

```go
// Has billing address
filter.User.Billing.IsNotNil()

// No billing address
filter.User.Billing.IsNil()
```

### Optional Struct Field Access

```go
// Filter by optional struct's field (if struct is not nil)
filter.User.Billing.City.Equal("Chicago")
```

## Sorting

Structs themselves are not sortable, but nested fields are:

```go
// Sort by nested field
query.Order(by.User.Address.City.Asc())

// Sort by nested numeric field
query.Order(by.Product.Dimensions.Width.Desc())

// Multiple nested sorts
query.Order(
    by.User.Address.Country.Asc(),
    by.User.Address.City.Asc(),
)
```

## Anonymous Embedding

Go's anonymous embedding is supported:

```go
type BaseInfo struct {
    CreatedBy string
    UpdatedBy string
}

type Document struct {
    som.Node

    Title    string
    BaseInfo // Anonymous embedding
}

// Access as direct fields
filter.Document.CreatedBy.Equal("admin")
```

## Deeply Nested Structs

Structs can contain other structs:

```go
type ContactInfo struct {
    Email   string
    Phone   string
    Address Address  // Nested struct
}

type Company struct {
    som.Node

    Name    string
    Contact ContactInfo
}

// Access deeply nested fields
filter.Company.Contact.Address.City.Equal("San Francisco")
filter.Company.Contact.Email.EndsWith("@company.com")
```

## Common Patterns

### Filter by Location

```go
// Users in New York
nyUsers, _ := client.UserRepo().Query().
    Where(filter.User.Address.City.Equal("New York")).
    All(ctx)
```

### Filter by Country

```go
// Users in USA
usUsers, _ := client.UserRepo().Query().
    Where(filter.User.Address.Country.Equal("USA")).
    All(ctx)
```

### Users with Billing Address

```go
// Users with separate billing
withBilling, _ := client.UserRepo().Query().
    Where(filter.User.Billing.IsNotNil()).
    All(ctx)
```

### Dimension Filtering

```go
// Large products
largeProducts, _ := client.ProductRepo().Query().
    Where(
        filter.Product.Dimensions.Width.GreaterThan(100),
        filter.Product.Dimensions.Height.GreaterThan(100),
    ).
    All(ctx)

// Calculate volume (simplified)
wideProducts, _ := client.ProductRepo().Query().
    Where(filter.Product.Dimensions.Width.GreaterThan(50)).
    Order(by.Product.Dimensions.Width.Desc()).
    All(ctx)
```

## Complete Example

```go
package main

import (
    "context"
    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/filter"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create user with address
    user := &model.User{
        Name: "Alice",
        Address: model.Address{
            Street:  "123 Main St",
            City:    "New York",
            Country: "USA",
            ZipCode: "10001",
        },
    }
    client.UserRepo().Create(ctx, user)

    // Find by city
    nyUsers, _ := client.UserRepo().Query().
        Where(filter.User.Address.City.Equal("New York")).
        All(ctx)

    // Find by country
    usUsers, _ := client.UserRepo().Query().
        Where(filter.User.Address.Country.Equal("USA")).
        All(ctx)

    // Find by zip code prefix
    manhattan, _ := client.UserRepo().Query().
        Where(filter.User.Address.ZipCode.StartsWith("100")).
        All(ctx)

    // Users with billing address
    withBilling, _ := client.UserRepo().Query().
        Where(filter.User.Billing.IsNotNil()).
        All(ctx)

    // Users with billing in different city
    differentBilling, _ := client.UserRepo().Query().
        Where(
            filter.User.Billing.IsNotNil(),
            filter.User.Billing.City.NotEqual(filter.User.Address.City),
        ).
        All(ctx)

    // Sort by city
    sorted, _ := client.UserRepo().Query().
        Order(by.User.Address.City.Asc()).
        All(ctx)

    // Products by dimensions
    products, _ := client.ProductRepo().Query().
        Where(
            filter.Product.Dimensions.Width.LessThan(100),
            filter.Product.Dimensions.Height.LessThan(100),
        ).
        Order(by.Product.Dimensions.Width.Desc()).
        All(ctx)
}
```

## Filter Reference Table

### Struct Access

| Operation | Description | Returns |
|-----------|-------------|---------|
| `<FieldName>` | Access nested field | Field's filter type |
| `IsNil()` | Struct is null (ptr) | Bool filter |
| `IsNotNil()` | Struct is not null (ptr) | Bool filter |

### Nested Field Operations

All filter operations available for the nested field's type:

```go
// String fields in struct
filter.User.Address.City.Equal("NYC")
filter.User.Address.City.Contains("York")
filter.User.Address.City.Lowercase().StartsWith("new")

// Numeric fields in struct
filter.Product.Dimensions.Width.GreaterThan(10)
filter.Product.Dimensions.Width.Mul(2).LessThan(100)

// Bool fields in struct
filter.Config.Settings.IsEnabled.True()
```

## Best Practices

### Keep Structs Simple

Embedded structs work best for simple value objects:

```go
// Good - simple value object
type Address struct {
    Street  string
    City    string
    Country string
}

// Consider Node instead for entities with identity
type Location struct {
    som.Node  // Has its own ID, can be referenced
    Name      string
    Latitude  float64
    Longitude float64
}
```

### Avoid Deep Nesting

Limit nesting depth for maintainability:

```go
// Reasonable
user.Address.City

// Getting complex
user.Contact.Address.Location.Coordinates.Latitude
```

### Use Pointers for Optional

Optional embedded structs should be pointers:

```go
type User struct {
    som.Node

    Address  Address   // Always required
    Billing  *Address  // Optional
    Shipping *Address  // Optional
}
```
