package basic

import (
	"context"
	"errors"
	"testing"
	"time"

	"som.test/gen/som"
	"som.test/model"
	"gotest.tools/v3/assert"
)

// validValidated returns a Validated model that passes every check.
func validValidated(name string) *model.Validated {
	return &model.Validated{
		Name: name,
		Contact: model.ContactAddress{
			City:  "Berlin",
			Email: "someone@example.com",
		},
		Code:    "ABCD",
		Ignored: "too long to be a code",
	}
}

func TestValidateOnCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	node := validValidated("")

	err := client.ValidatedRepo().Create(ctx, node)

	var validationErr *som.ValidationError
	assert.Assert(t, errors.As(err, &validationErr), "expected a validation error, got: %v", err)
	assert.Equal(t, "", validationErr.Path)
	assert.Equal(t, "name must not be empty", validationErr.Err.Error())
	assert.Equal(t, "", string(node.ID()), "the record must not have been written")

	// The sentinels are what a caller matches that only needs to tell a
	// rejected write from a failed one.
	assert.Assert(t, errors.Is(err, som.ErrValidation))
	assert.Assert(t, errors.Is(err, som.ErrInvalid))
	assert.Assert(t, !errors.Is(err, som.ErrAssert))
}

func TestValidateNestedPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	tests := []struct {
		name string
		node func(*model.Validated)
		path string
	}{
		{
			name: "nested struct",
			node: func(v *model.Validated) { v.Contact.City = "" },
			path: "Contact",
		},
		{
			name: "value below a nested struct",
			node: func(v *model.Validated) { v.Contact.Email = "not-an-email" },
			path: "Contact.Email",
		},
		{
			name: "optional struct",
			node: func(v *model.Validated) {
				v.ContactPtr = &model.ContactAddress{City: "Berlin", Email: "nope"}
			},
			path: "ContactPtr.Email",
		},
		{
			name: "slice element",
			node: func(v *model.Validated) {
				v.Contacts = []model.ContactAddress{
					{City: "Berlin", Email: "someone@example.com"},
					{City: "Hamburg", Email: "nope"},
				}
			},
			path: "Contacts[1].Email",
		},
		{
			name: "named scalar",
			node: func(v *model.Validated) { v.Code = "AB" },
			path: "Code",
		},
		{
			name: "scalar slice element",
			node: func(v *model.Validated) { v.Codes = []model.ProductCode{"ABCD", "XY"} },
			path: "Codes[1]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node := validValidated("some name")
			test.node(node)

			err := client.ValidatedRepo().Create(ctx, node)

			var validationErr *som.ValidationError
			assert.Assert(t, errors.As(err, &validationErr), "expected a validation error, got: %v", err)
			assert.Equal(t, test.path, validationErr.Path)
		})
	}
}

// TestValidateSkipsOptOut makes sure that a field tagged with `som:"novalidate"`
// is left alone, even though its type validates itself.
func TestValidateSkipsOptOut(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	node := validValidated("some name")

	assert.NilError(t, client.ValidatedRepo().Create(ctx, node))
	assert.Assert(t, node.ID() != "")
}

func TestValidateOnUpdate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	node := validValidated("some name")
	assert.NilError(t, client.ValidatedRepo().Create(ctx, node))

	node.Code = "AB"

	err := client.ValidatedRepo().Update(ctx, node)

	var validationErr *som.ValidationError
	assert.Assert(t, errors.As(err, &validationErr), "expected a validation error, got: %v", err)
	assert.Equal(t, "Code", validationErr.Path)

	// The record in the database must still hold the value it was created with.
	stored, exists, err := client.ValidatedRepo().Read(ctx, string(node.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, model.ProductCode("ABCD"), stored.Code)
}

func TestValidateOnInsert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	nodes := []*model.Validated{
		validValidated("first"),
		validValidated(""),
	}

	err := client.ValidatedRepo().Insert(ctx, nodes)

	var validationErr *som.ValidationError
	assert.Assert(t, errors.As(err, &validationErr), "expected a validation error, got: %v", err)

	count, err := client.ValidatedRepo().Query().Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, count, "no record of the batch must have been written")
}

// TestValidateOnSink covers the write-only repository, which validates as well.
func TestValidateOnSink(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.ValidatedLogRepo().Create(ctx, &model.ValidatedLog{Code: "AB"})

	var validationErr *som.ValidationError
	assert.Assert(t, errors.As(err, &validationErr), "expected a validation error, got: %v", err)
	assert.Equal(t, "Code", validationErr.Path)

	assert.NilError(t, client.ValidatedLogRepo().Create(ctx, &model.ValidatedLog{Code: "ABCD"}))
}

// TestValidateOnRelate covers the edge write path.
func TestValidateOnRelate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	from := validValidated("from")
	assert.NilError(t, client.ValidatedRepo().Create(ctx, from))

	to := &model.AllTypes{FieldString: "to", FieldMonth: time.January}
	assert.NilError(t, client.AllTypesRepo().Create(ctx, to))

	edge := &model.ValidatedEdge{Validated: *from, AllTypes: *to}

	err := client.ValidatedRepo().Relate().Edges().Create(ctx, edge)

	var validationErr *som.ValidationError
	assert.Assert(t, errors.As(err, &validationErr), "expected a validation error, got: %v", err)
	assert.Equal(t, "", validationErr.Path)
}

// TestValidateOnReadOnly makes sure a record that was stored before a rule
// existed can still be read, so that it can be fixed.
func TestValidateNotOnRead(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	node := validValidated("some name")
	assert.NilError(t, client.ValidatedRepo().Create(ctx, node))

	// The rule is bypassed by writing an invalid value directly.
	_, err := client.Raw(ctx,
		"UPDATE type::record('validated', $id) SET code = 'AB'",
		som.Params{"id": string(node.ID())},
	)
	assert.NilError(t, err)

	stored, exists, err := client.ValidatedRepo().Read(ctx, string(node.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, model.ProductCode("AB"), stored.Code)
}
