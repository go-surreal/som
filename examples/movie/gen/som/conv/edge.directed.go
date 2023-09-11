// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package conv

import (
	som "github.com/go-surreal/som"
	model "github.com/go-surreal/som/examples/movie/model"
)

type Directed struct {
	ID string `json:"id,omitempty"`
}

func FromDirected(data *model.Directed) *Directed {
	if data == nil {
		return nil
	}
	return &Directed{}
}

func ToDirected(data *Directed) *model.Directed {
	if data == nil {
		return nil
	}
	return &model.Directed{Edge: som.NewEdge(parseDatabaseID("directed", data.ID))}
}
