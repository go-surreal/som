module github.com/go-surreal/som

go 1.27.1

retract [v0.0.0, v0.19.99] // only the latest version is supported for now

require (
	github.com/dave/jennifer v1.7.1
	github.com/fxamacker/cbor/v2 v2.9.3
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/google/uuid v1.6.0
	github.com/iancoleman/strcase v0.3.0
	github.com/surrealdb/surrealdb.go v1.6.0
	github.com/urfave/cli/v3 v3.11.0
	github.com/wzshiming/gotype v0.8.0
	golang.org/x/exp v0.0.0-20260824195058-e88cd73687aa
	golang.org/x/mod v0.40.0
	gotest.tools/v3 v3.5.2
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
)
