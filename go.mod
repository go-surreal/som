module github.com/go-surreal/som

go 1.24.0

retract [v0.0.0, v0.9.99] // only the latest version is supported for now

require (
	github.com/dave/jennifer v1.7.1
	github.com/fxamacker/cbor/v2 v2.9.0
	github.com/google/uuid v1.6.0
	github.com/iancoleman/strcase v0.3.0
	github.com/surrealdb/surrealdb.go v1.0.0
	github.com/urfave/cli/v3 v3.6.1
	github.com/wzshiming/gotype v0.8.0
	golang.org/x/exp v0.0.0-20251113190631-e25ba8c21ef6
	golang.org/x/mod v0.30.0
	gotest.tools/v3 v3.5.2
)

require (
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/x448/float16 v0.8.4 // indirect
)
