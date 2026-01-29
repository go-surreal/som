package model

import (
	"net/url"

	"github.com/go-surreal/som/tests/basic/gen/som"
)

type SoftDeleteNode struct {
	som.Node
	som.Timestamps
	som.OptimisticLock
	som.SoftDelete

	Provider     string `som:"unique(provider_account)"`
	Account      string `som:"unique(provider_account)"`
	Name         string
	SomeURL      *url.URL
	SomeOtherURL url.URL
}

// SoftDeletePost is used to test fetch behavior with single-pointer soft-delete relations.
// Soft-delete filtering does NOT apply to fetched relations; all records are returned.
type SoftDeletePost struct {
	som.Node
	som.SoftDelete
	Title  string
	Author *SoftDeleteNode
}

// SoftDeleteBlogPost tests fetch behavior with slice relations to soft-delete models.
// Soft-delete filtering does NOT apply to fetched relations; all records are returned.
type SoftDeleteBlogPost struct {
	som.Node
	som.SoftDelete
	Title   string
	Authors []*SoftDeleteNode
}
