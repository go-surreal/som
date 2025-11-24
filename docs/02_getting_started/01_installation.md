# Installation

## Prerequisites

- Go 1.23 or later
- SurrealDB instance (local or remote)

## Install SOM

SOM is used as a code generator, so you run it via `go run`:

```bash
go run github.com/go-surreal/som/cmd/som@latest gen <input_dir> <output_dir>
```

Alternatively, install it as a binary:

```bash
go install github.com/go-surreal/som/cmd/som@latest
```

Then run:

```bash
som gen <input_dir> <output_dir>
```

## Add Generated Code to Your Project

After generation, add the output directory to your Go module. The generated code has its own `go.mod` file and can be imported like any other package.

## SurrealDB Setup

You need a running SurrealDB instance. For local development:

```bash
# Using Docker
docker run --rm -p 8000:8000 surrealdb/surrealdb:latest start --user root --pass root

# Or using the SurrealDB binary
surreal start --user root --pass root
```
