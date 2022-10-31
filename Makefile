# Makefile

gen:
	cd ./cmd/sdbgen && go run main.go ../../example/model ../../example/gen

run:
	cd ./cmd/sdbgen && go run main.go surreal

.PHONY: test
test:
	cd ./test && go run main.go
