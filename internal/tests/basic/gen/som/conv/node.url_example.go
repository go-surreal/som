// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package conv

import (
	v2 "github.com/fxamacker/cbor/v2"
	som "github.com/go-surreal/som"
	model "github.com/go-surreal/som/tests/basic/model"
)

type URLExample struct {
	ID           *som.ID `cbor:"id,omitempty"`
	SomeURL      *string `cbor:"some_url"`
	SomeOtherURL string  `cbor:"some_other_url"`
}

func FromURLExample(data model.URLExample) URLExample {
	return URLExample{
		SomeOtherURL: fromURL(data.SomeOtherURL),
		SomeURL:      fromURLPtr(data.SomeURL),
	}
}
func FromURLExamplePtr(data *model.URLExample) *URLExample {
	if data == nil {
		return nil
	}
	return &URLExample{
		SomeOtherURL: fromURL(data.SomeOtherURL),
		SomeURL:      fromURLPtr(data.SomeURL),
	}
}

func ToURLExample(data URLExample) model.URLExample {
	return model.URLExample{
		Node:         som.NewNode(data.ID),
		SomeOtherURL: toURL(data.SomeOtherURL),
		SomeURL:      toURLPtr(data.SomeURL),
	}
}
func ToURLExamplePtr(data *URLExample) *model.URLExample {
	if data == nil {
		return nil
	}
	return &model.URLExample{
		Node:         som.NewNode(data.ID),
		SomeOtherURL: toURL(data.SomeOtherURL),
		SomeURL:      toURLPtr(data.SomeURL),
	}
}

type urlexampleLink struct {
	URLExample
	ID *som.ID
}

func (f *urlexampleLink) MarshalCBOR() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return v2.Marshal(f.ID)
}

func (f *urlexampleLink) UnmarshalCBOR(data []byte) error {
	if err := v2.Unmarshal(data, &f.ID); err == nil {
		return nil
	}
	type alias urlexampleLink
	var link alias
	err := v2.Unmarshal(data, &link)
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
	return ToURLExample(res)
}

func fromURLExampleLinkPtr(link *urlexampleLink) *model.URLExample {
	if link == nil {
		return nil
	}
	res := URLExample(link.URLExample)
	out := ToURLExample(res)
	return &out
}

func toURLExampleLink(node model.URLExample) *urlexampleLink {
	if node.ID() == nil {
		return nil
	}
	link := urlexampleLink{URLExample: FromURLExample(node), ID: node.ID()}
	return &link
}

func toURLExampleLinkPtr(node *model.URLExample) *urlexampleLink {
	if node == nil || node.ID() == nil {
		return nil
	}
	link := urlexampleLink{URLExample: FromURLExample(*node), ID: node.ID()}
	return &link
}
