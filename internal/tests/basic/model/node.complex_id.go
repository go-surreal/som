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
	som.Node[WeatherKey]
	Temperature float64
}

type PersonKey struct {
	som.ObjectID
	Name string
	Age  int
}

type PersonObj struct {
	som.Node[PersonKey]
	Email string
}

type TeamMemberKey struct {
	som.ObjectID
	Member  AllTypes
	Forecast Weather
}

type TeamMember struct {
	som.Node[TeamMemberKey]
	Role string
}
