# Installation

## Prerequisites

- **Go 1.23 or later** - SOM uses generics extensively
- **SurrealDB 2.x** - Tested against version 2.3.10

## Install SOM Generator

### Option 1: Run Directly (Recommended)

Use `go run` to always get the latest version:

```bash
go run github.com/go-surreal/som/cmd/som@latest gen <input_dir> <output_dir>
```

### Option 2: Install Binary

Install globally for faster execution:

```bash
go install github.com/go-surreal/som/cmd/som@latest
```

Then run:

```bash
som gen <input_dir> <output_dir>
```

### Option 3: Pin Version

For reproducible builds, pin to a specific version:

```bash
go run github.com/go-surreal/som/cmd/som@v0.10.0 gen ./model ./gen/som
```

## CLI Flags

```bash
som gen [flags] <input_path> <output_path>

Flags:
  -v, --verbose    Show detailed generation progress
      --dry        Simulate generation without writing files
      --no-check   Skip version compatibility checks
```

## Project Setup

### 1. Create Model Directory

```bash
mkdir -p model
```

### 2. Define Your First Model

Create `model/user.go`:

```go
package model

import "github.com/go-surreal/som"

type User struct {
    som.Node

    Name  string
    Email string
}
```

### 3. Generate Code

```bash
go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som
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
docker run --rm -p 8000:8000 surrealdb/surrealdb:v2.3.10 \
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
    image: surrealdb/surrealdb:v2.3.10
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
