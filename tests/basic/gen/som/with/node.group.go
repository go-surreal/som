// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package with

import model "github.com/go-surreal/som/tests/basic/model"

var Group = group[model.Group]("")

type group[M any] string

func (n group[M]) fetch(M) {}
