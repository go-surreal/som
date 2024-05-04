// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package som

import "context"

type Client interface {
	URLExampleRepo() URLExampleRepo
	GroupRepo() GroupRepo
	FieldsLikeDBResponseRepo() FieldsLikeDBResponseRepo
	AllFieldTypesRepo() AllFieldTypesRepo
	ApplySchema(ctx context.Context) error
	Close()
}
