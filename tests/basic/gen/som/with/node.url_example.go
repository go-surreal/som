// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package with

import model "github.com/go-surreal/som/tests/basic/model"

var URLExample = urlexample[model.URLExample]("")

type urlexample[M any] string

func (n urlexample[M]) fetch(M) {}
