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

// SoftDeletePost is used to test fetch filtering on single-pointer soft-delete relations.
// Single relations require Go-side post-processing to filter soft-deleted records.
type SoftDeletePost struct {
	som.Node
	som.SoftDelete
	Title  string
	Author *SoftDeleteUser
}

// SoftDeleteBlogPost tests fetch filtering for slice relations to soft-delete models.
// Slice relations use DB-layer filtering (FETCH field[WHERE deleted_at IS NONE]).
type SoftDeleteBlogPost struct {
	som.Node
	som.SoftDelete
	Title   string
	Authors []*SoftDeleteUser
}
