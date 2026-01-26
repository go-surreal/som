package model

import "github.com/go-surreal/som/tests/basic/gen/som"

type SoftDeleteUser struct {
	som.Node
	som.SoftDelete
	Name string
}

type SoftDeleteComplete struct {
	som.Node
	som.Timestamps
	som.OptimisticLock
	som.SoftDelete
	Name string
}

// SoftDeletePost is used to test fetch behavior with single-pointer soft-delete relations.
// Soft-delete filtering does NOT apply to fetched relations; all records are returned.
type SoftDeletePost struct {
	som.Node
	som.SoftDelete
	Title  string
	Author *SoftDeleteUser
}

// SoftDeleteBlogPost tests fetch behavior with slice relations to soft-delete models.
// Soft-delete filtering does NOT apply to fetched relations; all records are returned.
type SoftDeleteBlogPost struct {
	som.Node
	som.SoftDelete
	Title   string
	Authors []*SoftDeleteUser
}
