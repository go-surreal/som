package model

import (
	"context"
	"errors"

	"github.com/go-surreal/som/tests/basic/gen/som"
)

type FieldsLikeDBResponse struct {
	som.Node

	Time   string
	Status string
	Detail string
	Result []string
}

type contextKey string

const AbortDeleteKey contextKey = "abortDelete"
const AfterDeleteCalledKey contextKey = "afterDeleteCalled"

func (f *FieldsLikeDBResponse) OnBeforeCreate(_ context.Context) error {
	f.Status = "[created]" + f.Status
	return nil
}

func (f *FieldsLikeDBResponse) OnAfterCreate(_ context.Context) error {
	f.Detail = f.Detail + "[after-create]"
	return nil
}

func (f *FieldsLikeDBResponse) OnBeforeUpdate(_ context.Context) error {
	f.Status = "[updated]" + f.Status
	return nil
}

func (f *FieldsLikeDBResponse) OnAfterUpdate(_ context.Context) error {
	f.Detail = f.Detail + "[after-update]"
	return nil
}

func (f *FieldsLikeDBResponse) OnBeforeDelete(ctx context.Context) error {
	if ptr, ok := ctx.Value(AbortDeleteKey).(*bool); ok && *ptr {
		return errors.New("delete aborted by model hook")
	}
	return nil
}

func (f *FieldsLikeDBResponse) OnAfterDelete(ctx context.Context) error {
	if ptr, ok := ctx.Value(AfterDeleteCalledKey).(*bool); ok {
		*ptr = true
	}
	return nil
}
