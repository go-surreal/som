# Numeric Type

Numeric types cover all Go integer and floating-point types with arithmetic operations, math functions, and type conversions.

## Overview

| Property | Value |
|----------|-------|
| Go Types | `int`, `int8`, `int16`, `int32`, `int64`, `uint8`, `uint16`, `uint32`, `float32`, `float64`, `rune` |
| Database Schema | `int` / `float` / `option<int>` / `option<float>` |
| CBOR Encoding | Direct |
| Sortable | Yes |

## Supported Types

### Integer Types

| Go Type | Range | DB Schema |
|---------|-------|-----------|
| `int` | Platform dependent | `int` |
| `int8` | -128 to 127 | `int` with assertion |
| `int16` | -32768 to 32767 | `int` with assertion |
| `int32` | -2B to 2B | `int` with assertion |
| `int64` | Full 64-bit | `int` |
| `rune` | Same as int32 | `int` with assertion |

### Unsigned Integer Types

| Go Type | Range | DB Schema |
|---------|-------|-----------|
| `uint8` / `byte` | 0 to 255 | `int` with assertion |
| `uint16` | 0 to 65535 | `int` with assertion |
| `uint32` | 0 to 4B | `int` with assertion |

> **Note**: `uint`, `uint64`, and `uintptr` are not supported due to SurrealDB integer limitations.

### Floating-Point Types

| Go Type | Precision | DB Schema |
|---------|-----------|-----------|
| `float32` | 32-bit | `float` |
| `float64` | 64-bit | `float` |

## Definition

```go
type Product struct {
    som.Node

    Price       float64  // Required float
    Quantity    int      // Required int
    Discount    *float32 // Optional float
    StockLevel  int16    // Range-checked int
}

type Metrics struct {
    som.Node

    Count8   int8
    Count16  int16
    Count32  int32
    Count64  int64
    Size8    uint8
    Size16   uint16
    Size32   uint32
    Value32  float32
    Value64  float64
}
```

## Schema

Generated SurrealDB schema:

```surql
DEFINE FIELD price ON product TYPE float;
DEFINE FIELD quantity ON product TYPE int;
DEFINE FIELD discount ON product TYPE option<float>;
DEFINE FIELD stock_level ON product TYPE int
    ASSERT $value >= -32768 AND $value <= 32767;
```

## Filter Operations

### Equality Operations

```go
// Exact match
filter.Product.Quantity.Equal(10)

// Not equal
filter.Product.Quantity.NotEqual(0)
```

### Set Membership

```go
// Value in set
filter.Product.Quantity.In(1, 5, 10, 25, 50)

// Value not in set
filter.Product.Quantity.NotIn(0, -1)
```

### Comparison Operations

```go
// Less than
filter.Product.Price.LessThan(100.0)

// Less than or equal
filter.Product.Price.LessThanEqual(99.99)

// Greater than
filter.Product.Quantity.GreaterThan(0)

// Greater than or equal
filter.Product.StockLevel.GreaterThanEqual(10)
```

### Arithmetic Operations

All arithmetic operations return a Float filter:

```go
// Addition
filter.Product.Price.Add(10.0).LessThan(100.0)

// Subtraction
filter.Product.Price.Sub(5.0).GreaterThan(50.0)

// Multiplication
filter.Product.Price.Mul(1.1).LessThan(110.0)  // 10% markup

// Division
filter.Product.Price.Div(2.0).LessThan(50.0)

// Exponentiation
filter.Metrics.Value.Raise(2).LessThan(100.0)  // Squared
```

### Math Functions

```go
// Absolute value
filter.Account.Balance.Abs().GreaterThan(100.0)

// Square root
filter.Metrics.Variance.Sqrt().LessThan(10.0)

// Ceiling (round up)
filter.Product.Price.Ceil().Equal(100.0)

// Floor (round down)
filter.Product.Price.Floor().Equal(99.0)

// Round to nearest
filter.Product.Price.Round().Equal(100.0)

// Fixed decimal places
filter.Product.Price.Fixed(2).Equal(99.99)
```

### Type Conversion

Convert between numeric types:

```go
// To integer types
filter.Product.Price.Int().Equal(99)
filter.Product.Price.Int8().GreaterThan(50)
filter.Product.Price.Int16().LessThan(1000)
filter.Product.Price.Int32().Equal(99)
filter.Product.Price.Int64().GreaterThan(0)

// To unsigned types
filter.Product.Quantity.Uint().GreaterThan(0)
filter.Product.Quantity.Uint8().LessThan(255)
filter.Product.Quantity.Uint16().GreaterThan(0)
filter.Product.Quantity.Uint32().LessThan(1000000)
filter.Product.Quantity.Uint64().GreaterThan(0)

// To float types
filter.Product.Quantity.Float32().GreaterThan(0.0)
filter.Product.Quantity.Float64().LessThan(100.0)
```

### Duration Conversion

Convert numbers to duration filters:

```go
// Interpret number as duration units
filter.Task.DelaySeconds.AsDurationSecs().LessThan(5 * time.Minute)
filter.Task.DelayMillis.AsDurationMillis().LessThan(time.Second)
filter.Config.TimeoutMins.AsDurationMins().GreaterThan(time.Hour)

// All duration conversions:
// AsDurationNanos() - nanoseconds
// AsDurationMicros() - microseconds
// AsDurationMillis() - milliseconds
// AsDurationSecs() - seconds
// AsDurationMins() - minutes
// AsDurationHours() - hours
// AsDurationDays() - days
// AsDurationWeeks() - weeks
```

### Time Conversion

Convert numbers to time filters:

```go
// Unix timestamp
filter.Event.Timestamp.AsTimeUnix().After(cutoffTime)

// Time from various units
filter.Event.EpochSecs.AsTimeSecs().Before(deadline)
filter.Event.EpochMillis.AsTimeMillis().After(startTime)
filter.Event.EpochMicros.AsTimeMicros().Before(endTime)
filter.Event.EpochNanos.AsTimeNanos().After(reference)
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
filter.Product.Discount.IsNil()

// Check if not nil
filter.Product.Discount.IsNotNil()
```

### Zero Value Check

```go
// Is zero
filter.Product.Quantity.Zero(true)

// Is not zero
filter.Product.Price.Zero(false)
```

## Sorting

```go
// Ascending (low to high)
query.Order(by.Product.Price.Asc())

// Descending (high to low)
query.Order(by.Product.Price.Desc())

// Multiple sort fields
query.Order(
    by.Product.InStock.Desc(),  // In-stock first
    by.Product.Price.Asc(),     // Then by price
)
```

## Method Chaining

Numeric filters support powerful chaining:

```go
// Calculate with tax and compare
filter.Product.Price.Mul(1.08).LessThan(100.0)

// Absolute difference from target
filter.Product.Price.Sub(targetPrice).Abs().LessThan(10.0)

// Round and compare
filter.Metrics.Average.Round().Equal(100)

// Complex calculation
filter.Product.Price.Mul(quantity).Sub(discount).LessThan(budget)
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

    // Find affordable products
    affordable, _ := client.ProductRepo().Query().
        Where(
            filter.Product.Price.LessThan(50.0),
            filter.Product.Quantity.GreaterThan(0),
        ).
        Order(by.Product.Price.Asc()).
        All(ctx)

    // Products with discount
    discounted, _ := client.ProductRepo().Query().
        Where(filter.Product.Discount.IsNotNil()).
        All(ctx)

    // Calculate final price with 10% tax
    underBudget, _ := client.ProductRepo().Query().
        Where(
            filter.Product.Price.Mul(1.1).LessThanEqual(100.0),
        ).
        All(ctx)

    // Find products near target price (within $5)
    targetPrice := 49.99
    nearTarget, _ := client.ProductRepo().Query().
        Where(
            filter.Product.Price.Sub(targetPrice).Abs().LessThan(5.0),
        ).
        All(ctx)

    // Low stock items
    lowStock, _ := client.ProductRepo().Query().
        Where(
            filter.Product.StockLevel.LessThan(10),
            filter.Product.StockLevel.GreaterThan(0),
        ).
        All(ctx)
}
```

## Filter Reference Table

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Equal(val)` | Exact match | Bool filter |
| `NotEqual(val)` | Not equal | Bool filter |
| `In(vals...)` | Value in set | Bool filter |
| `NotIn(vals...)` | Value not in set | Bool filter |
| `LessThan(val)` | Strictly less | Bool filter |
| `LessThanEqual(val)` | Less or equal | Bool filter |
| `GreaterThan(val)` | Strictly greater | Bool filter |
| `GreaterThanEqual(val)` | Greater or equal | Bool filter |
| `Add(val)` | Addition | Float filter |
| `Sub(val)` | Subtraction | Float filter |
| `Mul(val)` | Multiplication | Float filter |
| `Div(val)` | Division | Float filter |
| `Raise(val)` | Exponentiation | Float filter |
| `Abs()` | Absolute value | Float filter |
| `Sqrt()` | Square root | Float filter |
| `Ceil()` | Round up | Float filter |
| `Floor()` | Round down | Float filter |
| `Round()` | Round nearest | Float filter |
| `Fixed(places)` | Fixed decimals | Float filter |
| `Int()` | To int | Int filter |
| `Int8()` | To int8 | Int filter |
| `Int16()` | To int16 | Int filter |
| `Int32()` | To int32 | Int filter |
| `Int64()` | To int64 | Int filter |
| `Uint()` | To uint | Int filter |
| `Uint8()` | To uint8 | Int filter |
| `Uint16()` | To uint16 | Int filter |
| `Uint32()` | To uint32 | Int filter |
| `Uint64()` | To uint64 | Int filter |
| `Float32()` | To float32 | Float filter |
| `Float64()` | To float64 | Float filter |
| `AsDurationSecs()` | To duration | Duration filter |
| `AsDurationMillis()` | To duration | Duration filter |
| `AsDurationMins()` | To duration | Duration filter |
| `AsDurationHours()` | To duration | Duration filter |
| `AsDurationDays()` | To duration | Duration filter |
| `AsDurationWeeks()` | To duration | Duration filter |
| `AsDurationMicros()` | To duration | Duration filter |
| `AsDurationNanos()` | To duration | Duration filter |
| `AsTimeUnix()` | To time | Time filter |
| `AsTimeSecs()` | To time | Time filter |
| `AsTimeMillis()` | To time | Time filter |
| `AsTimeMicros()` | To time | Time filter |
| `AsTimeNanos()` | To time | Time filter |
| `Zero(bool)` | Check zero | Bool filter |
| `Truth()` | To boolean | Bool filter |
| `IsNil()` | Is null (ptr) | Bool filter |
| `IsNotNil()` | Not null (ptr) | Bool filter |
