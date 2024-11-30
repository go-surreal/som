module github.com/go-surreal/som

go 1.22.7

retract [v0.1.0, v0.6.2] // only the latest version is supported for now

require (
	github.com/dave/jennifer v1.7.0
	github.com/go-surreal/sdbc v0.6.1
	github.com/iancoleman/strcase v0.3.0
	github.com/urfave/cli/v2 v2.27.4
	github.com/wzshiming/gotype v0.7.3
	golang.org/x/exp v0.0.0-20240808152545-0cdaa3abc0fa
	golang.org/x/mod v0.20.0
	gotest.tools/v3 v3.5.1
)

require (
	github.com/coder/websocket v1.8.12 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/docker/docker v27.1.2+incompatible // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/sys v0.23.0 // indirect
)
