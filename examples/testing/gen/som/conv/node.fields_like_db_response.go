// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package conv

import (
	"encoding/json"
	som "github.com/marcbinz/som"
	model "github.com/marcbinz/som/examples/testing/model"
	"strings"
)

type FieldsLikeDBResponse struct {
	ID     string   `json:"id,omitempty"`
	Time   string   `json:"time"`
	Status string   `json:"status"`
	Detail string   `json:"detail"`
	Result []string `json:"result"`
}

func FromFieldsLikeDBResponse(data *model.FieldsLikeDBResponse) *FieldsLikeDBResponse {
	if data == nil {
		return nil
	}
	return &FieldsLikeDBResponse{
		Detail: data.Detail,
		Result: data.Result,
		Status: data.Status,
		Time:   data.Time,
	}
}

func ToFieldsLikeDBResponse(data *FieldsLikeDBResponse) *model.FieldsLikeDBResponse {
	if data == nil {
		return nil
	}
	return &model.FieldsLikeDBResponse{
		Detail: data.Detail,
		Node:   som.NewNode(parseDatabaseID("fields_like_db_response", data.ID)),
		Result: data.Result,
		Status: data.Status,
		Time:   data.Time,
	}
}

type fieldsLikeDBResponseLink struct {
	FieldsLikeDBResponse
	ID string
}

func (f *fieldsLikeDBResponseLink) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f.ID)
}

func (f *fieldsLikeDBResponseLink) UnmarshalJSON(data []byte) error {
	raw := string(data)
	if strings.HasPrefix(raw, "\"") && strings.HasSuffix(raw, "\"") {
		raw = raw[1 : len(raw)-1]
		f.ID = parseDatabaseID("fields_like_db_response", raw)
		return nil
	}
	type alias fieldsLikeDBResponseLink
	var link alias
	err := json.Unmarshal(data, &link)
	if err == nil {
		*f = fieldsLikeDBResponseLink(link)
	}
	return err
}

func fromFieldsLikeDBResponseLink(link *fieldsLikeDBResponseLink) model.FieldsLikeDBResponse {
	if link == nil {
		return model.FieldsLikeDBResponse{}
	}
	res := FieldsLikeDBResponse(link.FieldsLikeDBResponse)
	out := ToFieldsLikeDBResponse(&res)
	return *out
}

func fromFieldsLikeDBResponseLinkPtr(link *fieldsLikeDBResponseLink) *model.FieldsLikeDBResponse {
	if link == nil {
		return nil
	}
	res := FieldsLikeDBResponse(link.FieldsLikeDBResponse)
	return ToFieldsLikeDBResponse(&res)
}

func toFieldsLikeDBResponseLink(node model.FieldsLikeDBResponse) *fieldsLikeDBResponseLink {
	if node.ID() == "" {
		return nil
	}
	link := fieldsLikeDBResponseLink{FieldsLikeDBResponse: *FromFieldsLikeDBResponse(&node), ID: buildDatabaseID("fields_like_db_response", node.ID())}
	return &link
}

func toFieldsLikeDBResponseLinkPtr(node *model.FieldsLikeDBResponse) *fieldsLikeDBResponseLink {
	if node == nil || node.ID() == "" {
		return nil
	}
	link := fieldsLikeDBResponseLink{FieldsLikeDBResponse: *FromFieldsLikeDBResponse(node), ID: buildDatabaseID("fields_like_db_response", node.ID())}
	return &link
}
