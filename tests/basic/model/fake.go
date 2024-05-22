package model

import "github.com/go-surreal/som"

type Fake struct {
	som.Node

	Name string `som:"type=name"`
}

type FakeFaker struct{}

func NewFake()
