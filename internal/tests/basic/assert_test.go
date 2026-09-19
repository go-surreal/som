package basic

import (
	"context"
	"errors"
	"testing"

	"gotest.tools/v3/assert"
	"som.test/gen/som"
	"som.test/model"
)

// validConstrained returns a record satisfying every constraint, so that a test
// can violate exactly one of them.
func validConstrained() model.Constrained {
	return model.Constrained{
		Name:  "Valid Name",
		Code:  "ABCD",
		Age:   42,
		Ratio: 1.5,
		Tags:  []string{"one", "two"},
		Phone: "+49-123",
		Title: "A Title",
	}
}

func TestAssertAcceptsValidRecord(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := validConstrained()

	err := client.ConstrainedRepo().Create(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, record.ID() != "")
}

func TestAssertReportsViolatedField(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	tests := []struct {
		name    string
		mutate  func(record *model.Constrained)
		field   string
		message string
	}{
		{
			name: "length below minimum",
			mutate: func(record *model.Constrained) {
				record.Name = "ab"
			},
			field:   "name",
			message: "length must be between 3 and 64",
		},
		{
			name: "length not exact",
			mutate: func(record *model.Constrained) {
				record.Code = "ABC"
			},
			field:   "code",
			message: "length must be exactly 4",
		},
		{
			name: "number above maximum",
			mutate: func(record *model.Constrained) {
				record.Age = 131
			},
			field:   "age",
			message: "must be between 0 and 130",
		},
		{
			name: "number below minimum",
			mutate: func(record *model.Constrained) {
				record.Ratio = -0.5
			},
			field:   "ratio",
			message: "must be at least 0",
		},
		{
			name: "too many items",
			mutate: func(record *model.Constrained) {
				record.Tags = []string{"1", "2", "3", "4", "5", "6"}
			},
			field:   "tags",
			message: "length must be at most 5",
		},
		{
			name: "pattern mismatch",
			mutate: func(record *model.Constrained) {
				record.Phone = "not a phone"
			},
			field:   "phone",
			message: "must only contain digits, plus and minus",
		},
		{
			name: "raw expression",
			mutate: func(record *model.Constrained) {
				record.Title = "TBD"
			},
			field:   "title",
			message: "must conform to no_placeholder",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			record := validConstrained()
			test.mutate(&record)

			err := client.ConstrainedRepo().Create(ctx, &record)
			assert.Assert(t, err != nil, "expected the write to be rejected")
			assert.Assert(t, errors.Is(err, som.ErrAssert), "expected ErrAssert, got: %v", err)
			assert.Assert(t, errors.Is(err, som.ErrInvalid), "expected ErrInvalid, got: %v", err)

			var assertErr *som.AssertError
			assert.Assert(t, errors.As(err, &assertErr), "expected AssertError, got: %v", err)
			assert.Equal(t, test.field, assertErr.Field)
			assert.Equal(t, test.message, assertErr.Message)
		})
	}
}

// An optional field carries its constraint only once it holds a value.
func TestAssertSkipsAbsentOptionalValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := validConstrained()
	record.Nickname = nil
	record.Tags = nil

	err := client.ConstrainedRepo().Create(ctx, &record)
	assert.NilError(t, err)

	tooShort := "x"
	record.Nickname = &tooShort

	err = client.ConstrainedRepo().Update(ctx, &record)
	assert.Assert(t, errors.Is(err, som.ErrAssert), "expected ErrAssert, got: %v", err)
}

// The constraints declared via tags compose with the ones som derives from the
// Go type, rather than replacing them.
func TestAssertComposesWithBuiltinRange(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := validConstrained()
	record.Age = -1

	err := client.ConstrainedRepo().Create(ctx, &record)

	var assertErr *som.AssertError
	assert.Assert(t, errors.As(err, &assertErr), "expected AssertError, got: %v", err)
	assert.Equal(t, "age", assertErr.Field)
}
