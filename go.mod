module github.com/go-surreal/som

go 1.23.7

retract [v0.1.0, v0.6.4] // only the latest version is supported for now

require (
	github.com/dave/jennifer v1.7.1
	github.com/go-surreal/sdbc v0.9.2
	github.com/iancoleman/strcase v0.3.0
	github.com/urfave/cli/v2 v2.27.6
	github.com/wzshiming/gotype v0.7.4
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0
	golang.org/x/mod v0.24.0
	gotest.tools/v3 v3.5.2
)

require (
	github.com/coder/websocket v1.8.13 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/fxamacker/cbor/v2 v2.8.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)
