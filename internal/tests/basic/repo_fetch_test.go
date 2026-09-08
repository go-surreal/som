package basic

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
	"som.test/gen/som"
	"som.test/gen/som/with"
	"som.test/model"
)

func TestRepoFetchLink(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	author := model.SpecialTypes{Name: "Author"}
	err := client.SpecialTypesRepo().Create(ctx, &author)
	assert.NilError(t, err)

	post := model.SpecialRelation{Title: "Post", Author: &author}
	err = client.SpecialRelationRepo().Create(ctx, &post)
	assert.NilError(t, err)

	read, exists, err := client.SpecialRelationRepo().Read(ctx, string(post.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Assert(t, read.Author.IsPartial())

	err = client.SpecialRelationRepo().Fetch(ctx, read, with.SpecialRelation.Author())
	assert.NilError(t, err)

	assert.Assert(t, !read.Author.IsPartial(), "a fetched link is fully loaded")
	assert.Assert(t, read.Author.Marker().Has(som.MarkerLoaded))
	assert.Equal(t, read.Author.Name, author.Name)
}

func TestRepoFetchWithoutFields(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	author := model.SpecialTypes{Name: "Author"}
	err := client.SpecialTypesRepo().Create(ctx, &author)
	assert.NilError(t, err)

	post := model.SpecialRelation{Title: "Post", Author: &author}
	err = client.SpecialRelationRepo().Create(ctx, &post)
	assert.NilError(t, err)

	read, exists, err := client.SpecialRelationRepo().Read(ctx, string(post.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)

	err = client.SpecialRelationRepo().Fetch(ctx, read)
	assert.NilError(t, err)
	assert.Assert(t, read.Author.IsPartial(), "without fetch fields nothing is resolved")

	assert.Error(t,
		client.SpecialRelationRepo().Fetch(ctx, nil, with.SpecialRelation.Author()),
		"the passed node must not be nil",
	)
	assert.Error(t,
		client.SpecialRelationRepo().Fetch(ctx, &model.SpecialRelation{}, with.SpecialRelation.Author()),
		"cannot fetch SpecialRelation without existing record ID",
	)
}
