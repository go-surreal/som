# Byte Type

The byte type handles single byte values and byte slices for binary data storage.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `byte` / `*byte` or `[]byte` / `*[]byte` |
| Database Schema | `int` / `option<int>` (single) or `bytes` (slice) |
| CBOR Encoding | Direct |
| Sortable | Yes (single byte) / No (slice) |

## Single Byte

### Definition

```go
type Packet struct {
    som.Node

    Header    byte   // Required single byte
    TypeFlag  *byte  // Optional single byte
}
```

### Schema

Generated SurrealDB schema with range validation:

```surql
DEFINE FIELD header ON packet TYPE int ASSERT $value >= 0 AND $value <= 255;
DEFINE FIELD type_flag ON packet TYPE option<int>
    ASSERT $value == NONE OR $value == NULL OR ($value >= 0 AND $value <= 255);
```

### Filter Operations

Single byte fields support base filter operations:

```go
// Equality
where.Packet.Header.Equal(0xFF)
where.Packet.Header.NotEqual(0x00)

// Set membership
where.Packet.Header.In([]byte{0x01, 0x02, 0x03})
where.Packet.Header.NotIn([]byte{0x00})

// Zero check
where.Packet.Header.Zero(true)   // Is 0x00
where.Packet.Header.Zero(false)  // Is not 0x00
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
where.Packet.TypeFlag.IsNil()

// Check if not nil
where.Packet.TypeFlag.IsNotNil()
```

### Sorting

```go
// Ascending
query.Order(by.Packet.Header.Asc())

// Descending
query.Order(by.Packet.Header.Desc())
```

## Byte Slice

### Definition

```go
type Document struct {
    som.Node

    Data     []byte   // Required binary data
    Checksum *[]byte  // Optional binary data
}
```

### Schema

Byte slices are stored as the SurrealDB `bytes` type:

```surql
DEFINE FIELD data ON document TYPE option<bytes>;
DEFINE FIELD checksum ON document TYPE option<bytes>;
```

### Filter Operations

#### Base Operations

```go
// Equality
where.Document.Data.Equal([]byte{0x01, 0x02, 0x03})
where.Document.Data.NotEqual([]byte{})

// Set membership
where.Document.Data.In([][]byte{data1, data2})
where.Document.Data.NotIn([][]byte{invalidData})

// Zero check
where.Document.Data.Zero(true)   // Is empty/nil
where.Document.Data.Zero(false)  // Has data
```

#### Base64 Encoding

Convert byte slice to base64 string for string operations:

```go
// Encode to base64 and compare
where.Document.Data.Base64Encode().Equal("SGVsbG8gV29ybGQ=")

// Base64 with string operations
where.Document.Data.Base64Encode().StartsWith("SGVs")
where.Document.Data.Base64Encode().Contains("bG8=")
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
where.Document.Checksum.IsNil()

// Check if not nil
where.Document.Checksum.IsNotNil()
```

## Creating Byte Values

```go
// Single byte
packet := &model.Packet{
    Header: 0x7F,
}

// Optional single byte
flag := byte(0x01)
packet.TypeFlag = &flag

// Byte slice
document := &model.Document{
    Data: []byte("Hello World"),
}

// From hex
data, _ := hex.DecodeString("48656c6c6f")
document.Data = data
```

## Common Patterns

### Filter by Header Value

```go
// Packets with specific header
packets, _ := client.PacketRepo().Query().
    Filter(where.Packet.Header.Equal(0x01)).
    All(ctx)
```

### Filter by Data Content

```go
// Documents with specific data
docs, _ := client.DocumentRepo().Query().
    Filter(where.Document.Data.Equal(expectedData)).
    All(ctx)
```

### Base64 String Matching

```go
// Find by base64 prefix
docs, _ := client.DocumentRepo().Query().
    Filter(where.Document.Data.Base64Encode().StartsWith("SGVs")).
    All(ctx)
```

### Documents with Checksum

```go
// Documents that have a checksum
withChecksum, _ := client.DocumentRepo().Query().
    Filter(where.Document.Checksum.IsNotNil()).
    All(ctx)
```

## Complete Example

```go
package main

import (
    "context"
    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/where"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create document with binary data
    doc := &model.Document{
        Data: []byte("Hello World"),
    }
    client.DocumentRepo().Create(ctx, doc)

    // Find by exact data
    found, exists, _ := client.DocumentRepo().Query().
        Filter(where.Document.Data.Equal([]byte("Hello World"))).
        First(ctx)

    // Find by base64 encoding
    base64Docs, _ := client.DocumentRepo().Query().
        Filter(where.Document.Data.Base64Encode().Equal("SGVsbG8gV29ybGQ=")).
        All(ctx)

    // Find documents with checksum
    withChecksum, _ := client.DocumentRepo().Query().
        Filter(where.Document.Checksum.IsNotNil()).
        All(ctx)

    // Create packet with single byte
    packet := &model.Packet{
        Header: 0x01,
    }
    client.PacketRepo().Create(ctx, packet)

    // Filter by header value
    packets, _ := client.PacketRepo().Query().
        Filter(where.Packet.Header.Equal(0x01)).
        Order(by.Packet.Header.Asc()).
        All(ctx)
}
```

## Filter Reference Table

### Single Byte

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Equal(val)` | Exact match | Bool filter |
| `NotEqual(val)` | Not equal | Bool filter |
| `In(vals...)` | Value in set | Bool filter |
| `NotIn(vals...)` | Value not in set | Bool filter |
| `Zero(bool)` | Check if zero | Bool filter |
| `Truth()` | To boolean | Bool filter |
| `IsNil()` | Is null (ptr only) | Bool filter |
| `IsNotNil()` | Not null (ptr only) | Bool filter |

### Byte Slice

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Equal(val)` | Exact match | Bool filter |
| `NotEqual(val)` | Not equal | Bool filter |
| `In(vals...)` | Value in set | Bool filter |
| `NotIn(vals...)` | Value not in set | Bool filter |
| `Zero(bool)` | Check if empty | Bool filter |
| `Truth()` | To boolean | Bool filter |
| `Base64Encode()` | Convert to base64 | String filter |
| `IsNil()` | Is null (ptr only) | Bool filter |
| `IsNotNil()` | Not null (ptr only) | Bool filter |
