module github.com/go-surreal/som/tests/basic

go 1.22

toolchain go1.22.6

replace github.com/go-surreal/som => ../../../

require (
	github.com/docker/docker v27.1.2+incompatible
	github.com/fxamacker/cbor/v2 v2.7.0
	github.com/go-surreal/sdbc v0.6.1
	github.com/go-surreal/som v0.6.2
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.6.0
	github.com/testcontainers/testcontainers-go v0.32.0
	gotest.tools/v3 v3.5.1
)

require (
	github.com/coder/websocket v1.8.12 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/exp v0.0.0-20240808152545-0cdaa3abc0fa // indirect
	golang.org/x/sys v0.23.0 // indirect
)
