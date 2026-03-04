package basic

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"gotest.tools/v3/assert"
)

func TestServerError_TypeAlias(t *testing.T) {
	t.Parallel()

	se := som.ServerError{
		Code:    -32000,
		Message: "An error occurred",
		Kind:    "Thrown",
	}

	assert.Equal(t, "An error occurred", se.Error())
	assert.Equal(t, "Thrown", se.Kind)
	assert.Equal(t, -32000, se.Code)
}

func TestServerError_ErrorsAs(t *testing.T) {
	t.Parallel()

	se := som.ServerError{
		Code:    -32000,
		Message: "An error occurred",
		Kind:    "Thrown",
	}
	wrapped := fmt.Errorf("operation failed: %w", se)

	var extracted som.ServerError
	assert.Assert(t, errors.As(wrapped, &extracted), "errors.As should extract ServerError from wrapped error")
	assert.Equal(t, "An error occurred", extracted.Message)
	assert.Equal(t, "Thrown", extracted.Kind)
}

func TestServerError_CauseChain(t *testing.T) {
	t.Parallel()

	se := som.ServerError{
		Code:    -32000,
		Message: "outer error",
		Kind:    "QueryError",
		Cause: &som.ServerError{
			Code:    -32000,
			Message: "inner: optimistic_lock_failed",
			Kind:    "Thrown",
		},
	}

	assert.Equal(t, "outer error: inner: optimistic_lock_failed", se.Error())
	assert.Equal(t, "inner: optimistic_lock_failed", se.Cause.Message)
}

func TestServerError_ErrorsIs(t *testing.T) {
	t.Parallel()

	se := som.ServerError{
		Code:    -32000,
		Message: "some error",
		Kind:    "NotFound",
	}
	wrapped := fmt.Errorf("operation failed: %w", se)

	assert.Assert(t, errors.Is(wrapped, som.ServerError{}),
		"errors.Is should match any ServerError in the chain")
}

func TestServerError_NotPresent(t *testing.T) {
	t.Parallel()

	plainErr := errors.New("plain error")

	var se som.ServerError
	assert.Assert(t, !errors.As(plainErr, &se),
		"errors.As should not match when no ServerError in chain")
}
