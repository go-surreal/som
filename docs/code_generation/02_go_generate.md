# go:generate Integration

Integrate SOM code generation into your Go build process using `go:generate`.

## Setup

Add a generate directive to a file in your project (commonly `generate.go` or in your model package):

```go
//go:generate go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som

package main
```

## Running Generation

Execute all generate directives in your project:

```bash
go generate ./...
```

Or run from a specific directory:

```bash
go generate ./model/...
```

## Placement Options

### In root directory

Create a `generate.go` file:

```go
//go:generate go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som

package main
```

### In model package

Add the directive to any `.go` file in your model directory:

```go
//go:generate go run github.com/go-surreal/som/cmd/som@latest gen . ../gen/som

package model
```

## CI/CD Integration

Add generation to your CI pipeline to verify generated code is up to date:

```yaml
# .github/workflows/ci.yml
steps:
  - uses: actions/checkout@v4
  - uses: actions/setup-go@v5
    with:
      go-version: '1.23'
  - run: go generate ./...
  - run: git diff --exit-code  # Fail if generated code differs
```

## Multiple Model Directories

If you have multiple model directories, add a directive for each:

```go
//go:generate go run github.com/go-surreal/som/cmd/som@latest gen ./models/users ./gen/users
//go:generate go run github.com/go-surreal/som/cmd/som@latest gen ./models/products ./gen/products

package main
```
