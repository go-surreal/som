module github.com/go-surreal/som

go 1.23.7

retract [v0.1.0, v0.7.0] // only the latest version is supported for now

require (
	github.com/dave/jennifer v1.7.1
	github.com/fxamacker/cbor/v2 v2.9.0
	github.com/go-surreal/sdbc v0.9.4
	github.com/google/uuid v1.6.0
	github.com/iancoleman/strcase v0.3.0
	github.com/urfave/cli/v3 v3.6.0
	github.com/wzshiming/gotype v0.8.0
	golang.org/x/exp a4bb9ffd2546
	golang.org/x/mod v0.29.0
	gotest.tools/v3 v3.5.2
)

require (
	github.com/coder/websocket v1.8.13 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
)
