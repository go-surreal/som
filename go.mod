module github.com/go-surreal/som

go 1.23.7

retract [v0.1.0, v0.7.0] // only the latest version is supported for now

require (
	github.com/dave/jennifer v1.7.1
	github.com/fxamacker/cbor/v2 v2.8.0
	github.com/go-surreal/sdbc v0.9.3
	github.com/google/uuid v1.6.0
	github.com/iancoleman/strcase v0.3.0
	github.com/urfave/cli/v2 v2.27.6
	github.com/wzshiming/gotype v0.7.4
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6
	golang.org/x/mod v0.24.0
	gotest.tools/v3 v3.5.2
)

require (
	github.com/coder/websocket v1.8.13 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
)
