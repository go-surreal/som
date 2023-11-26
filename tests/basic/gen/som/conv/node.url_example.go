// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package conv

import (
	"encoding/json"
	som "github.com/go-surreal/som"
	model "github.com/go-surreal/som/tests/basic/model"
	"strings"
)

type URLExample struct {
	ID           string  `json:"id,omitempty"`
	SomeURL      *string `json:"some_url"`
	SomeOtherURL string  `json:"some_other_url"`
}

func FromURLExample(data *model.URLExample) *URLExample {
	if data == nil {
		return nil
	}
	return &URLExample{
		SomeOtherURL: data.SomeOtherURL.String(),
		SomeURL:      urlPtr(data.SomeURL),
	}
}

func ToURLExample(data *URLExample) *model.URLExample {
	if data == nil {
		return nil
	}
	return &model.URLExample{
		Node:         som.NewNode(parseDatabaseID("url_example", data.ID)),
		SomeOtherURL: parseURL(data.SomeOtherURL),
		SomeURL:      ptrFunc(parseURL)(data.SomeURL),
	}
}

type urlexampleLink struct {
	URLExample
	ID string
}

func (f *urlexampleLink) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f.ID)
}

func (f *urlexampleLink) UnmarshalJSON(data []byte) error {
	raw := string(data)
	if strings.HasPrefix(raw, "\"") && strings.HasSuffix(raw, "\"") {
		raw = raw[1 : len(raw)-1]
		f.ID = parseDatabaseID("url_example", raw)
		return nil
	}
	type alias urlexampleLink
	var link alias
	err := json.Unmarshal(data, &link)
	if err == nil {
		*f = urlexampleLink(link)
	}
	return err
}

func fromURLExampleLink(link *urlexampleLink) model.URLExample {
	if link == nil {
		return model.URLExample{}
	}
	res := URLExample(link.URLExample)
	out := ToURLExample(&res)
	return *out
}

func fromURLExampleLinkPtr(link *urlexampleLink) *model.URLExample {
	if link == nil {
		return nil
	}
	res := URLExample(link.URLExample)
	return ToURLExample(&res)
}

func toURLExampleLink(node model.URLExample) *urlexampleLink {
	if node.ID() == "" {
		return nil
	}
	link := urlexampleLink{URLExample: *FromURLExample(&node), ID: buildDatabaseID("url_example", node.ID())}
	return &link
}

func toURLExampleLinkPtr(node *model.URLExample) *urlexampleLink {
	if node == nil || node.ID() == "" {
		return nil
	}
	link := urlexampleLink{URLExample: *FromURLExample(node), ID: buildDatabaseID("url_example", node.ID())}
	return &link
}