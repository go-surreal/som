# Makefile

gen:
	cd ./cmd/sdbgen && go run main.go ../../example/model ../../example/gen
