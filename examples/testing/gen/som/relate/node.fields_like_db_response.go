// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package relate

func NewFieldsLikeDBResponse(db Database) *FieldsLikeDBResponse {
	return &FieldsLikeDBResponse{db: db}
}

type FieldsLikeDBResponse struct {
	db Database
}