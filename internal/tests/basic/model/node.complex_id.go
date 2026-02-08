package model

import (
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som"
)

type WeatherKey struct {
	som.ArrayID
	City string
	Date time.Time
}

type Weather struct {
	som.CustomNode[WeatherKey]
	Temperature float64
}

type PersonKey struct {
	som.ObjectID
	Name string
	Age  int
}

type PersonObj struct {
	som.CustomNode[PersonKey]
	Email string
}
