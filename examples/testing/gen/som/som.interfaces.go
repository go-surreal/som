// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package som

type Client interface {
	URLExampleRepo() URLExampleRepo
	FieldsLikeDBResponseRepo() FieldsLikeDBResponseRepo
	ApplySchema() error
	Close()
}
