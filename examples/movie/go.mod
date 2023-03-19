module github.com/marcbinz/som/examples/movie

go 1.19

replace github.com/marcbinz/som => ../../

replace github.com/docker/docker => github.com/docker/docker v20.10.3-0.20221021173910-5aac513617f0+incompatible // 22.06 branch

require (
	github.com/google/uuid v1.3.0
	github.com/marcbinz/som v0.0.0
	github.com/surrealdb/surrealdb.go v0.1.2-0.20230309175011-e0397e60fdc8
)

require github.com/gorilla/websocket v1.5.0 // indirect
