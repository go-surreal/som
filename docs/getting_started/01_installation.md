# Installation

## Prerequisites

- **Go 1.25 or later** - SOM uses generics and iterators extensively
- **SurrealDB 3.x** - Tested against version 3.0.0

## Install SOM Generator

### Option 1: Run Directly (Recommended)

Use `go run` to always get the latest version:

```bash
go run github.com/go-surreal/som@latest -i <input_dir>
```

### Option 2: Install Binary

Install globally for faster execution:

```bash
go install github.com/go-surreal/som@latest
```

Then run:

```bash
som -i <input_dir>
```

### Option 3: Pin Version

For reproducible builds, pin to a specific version:

```bash
go run github.com/go-surreal/som@v0.16.0 -i ./model
```

## CLI Flags

```bash
som [flags]

Flags:
  -i, --in             Input model directory (relative to go.mod)
  -o, --out            Output directory (default: gen/som)
  -v, --verbose        Show detailed generation progress
      --dry            Simulate generation without writing files
      --no-check       Skip version compatibility checks
      --no-count-index Disable automatic COUNT index generation
      --wire           Override wire generation: "no", "google", or "goforj"
```

The tool auto-detects the closest `go.mod` file and resolves paths relative to it. If `--in` is omitted, only static base files are generated (initialization mode).

## Project Setup

### 1. Create Model Directory

```bash
mkdir -p model
```

### 2. Define Your First Model

Create `model/user.go`:

```go
package model

import "yourproject/gen/som"

type User struct {
    som.Node[som.ULID]

    Name  string
    Email string
}
```

### 3. Generate Code

First, generate the static base files:

```bash
go run github.com/go-surreal/som@latest
```

Then generate model-specific code:

```bash
go run github.com/go-surreal/som@latest -i ./model
```

### 4. Import Generated Code

The generated code is a self-contained Go module. Import it in your application:

```go
import (
    "yourproject/gen/som"
    "yourproject/gen/som/filter"
    "yourproject/gen/som/by"
)
```

## SurrealDB Setup

### Using Docker (Recommended)

```bash
docker run --rm -p 8000:8000 surrealdb/surrealdb:v3.0.0 \
    start --user root --pass root
```

### Using Binary

Download from [surrealdb.com/install](https://surrealdb.com/install):

```bash
surreal start --user root --pass root --bind 0.0.0.0:8000
```

### Using Docker Compose

```yaml
version: '3'
services:
  surrealdb:
    image: surrealdb/surrealdb:v3.0.0
    ports:
      - "8000:8000"
    command: start --user root --pass root
```

## Verify Installation

Test that everything works:

```go
package main

import (
    "context"
    "log"

    "yourproject/gen/som"
)

func main() {
    ctx := context.Background()

    client, err := som.NewClient(ctx, som.Config{
        Address:   "ws://localhost:8000",
        Username:  "root",
        Password:  "root",
        Namespace: "test",
        Database:  "test",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    log.Println("Connected successfully!")
}
```

## Next Steps

- Follow the [Quick Start](02_quick_start.md) guide
- Learn about [Core Concepts](03_concepts.md)
