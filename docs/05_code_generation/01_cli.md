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
├── go.mod              # Module file for generated code
├── client.go           # Database client
├── where/              # Filter helpers
├── order/              # Ordering helpers
├── field/              # Field accessors
└── internal/           # Internal utilities
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
