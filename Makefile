# Makefile

gen:
	cd ./cmd/somgen && go run main.go ../../example/model ../../example/gen/som

run:
	cd ./cmd/somgen && go run main.go surreal

.PHONY: test
test:
	cd ./test && go run main.go
