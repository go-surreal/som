package basic

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
	"som.test/gen/som"
	"som.test/gen/som/with"
	"som.test/model"
)

func TestRepoResolveLink(t *testing.T) {
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

	err = client.SpecialRelationRepo().Resolve(ctx, read, with.SpecialRelation.Author())
	assert.NilError(t, err)

	assert.Assert(t, !read.Author.IsPartial(), "a resolved link is fully loaded")
	assert.Assert(t, read.Author.Marker().Has(som.MarkerLoaded))
	assert.Equal(t, read.Author.Name, author.Name)
}

func TestRepoResolveKeepsEarlierRelations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	author := model.SpecialTypes{Name: "Author"}
	err := client.SpecialTypesRepo().Create(ctx, &author)
	assert.NilError(t, err)

	reviewer := model.SpecialTypes{Name: "Reviewer"}
	err = client.SpecialTypesRepo().Create(ctx, &reviewer)
	assert.NilError(t, err)

	post := model.SpecialRelation{
		Title:   "Post",
		Author:  &author,
		Authors: []*model.SpecialTypes{&reviewer},
	}
	err = client.SpecialRelationRepo().Create(ctx, &post)
	assert.NilError(t, err)

	read, exists, err := client.SpecialRelationRepo().Read(ctx, string(post.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)

	err = client.SpecialRelationRepo().Resolve(ctx, read, with.SpecialRelation.Author())
	assert.NilError(t, err)
	assert.Assert(t, !read.Author.IsPartial())

	// Resolving another relation must not drop the one resolved before, even
	// though it is not part of the second response.
	err = client.SpecialRelationRepo().Resolve(ctx, read, with.SpecialRelation.Authors())
	assert.NilError(t, err)

	assert.Assert(t, !read.Authors[0].IsPartial())
	assert.Equal(t, read.Authors[0].Name, reviewer.Name)
	assert.Assert(t, !read.Author.IsPartial(), "an earlier resolved relation stays resolved")
	assert.Equal(t, read.Author.Name, author.Name)

	// The other fields of the record must survive the merge as well.
	assert.Equal(t, read.Title, post.Title)
	assert.Equal(t, read.ID(), post.ID())
}

func TestRepoResolveIsIdempotent(t *testing.T) {
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

	err = client.SpecialRelationRepo().Resolve(ctx, read, with.SpecialRelation.Author())
	assert.NilError(t, err)

	// The relation is loaded, so the second call must not touch the database.
	// Renaming the record remotely makes a re-read visible.
	author.Name = "Renamed"
	err = client.SpecialTypesRepo().Update(ctx, &author)
	assert.NilError(t, err)

	err = client.SpecialRelationRepo().Resolve(ctx, read, with.SpecialRelation.Author())
	assert.NilError(t, err)
	assert.Equal(t, read.Author.Name, "Author", "a resolved relation is not read again")

	err = client.SpecialTypesRepo().Refresh(ctx, read.Author)
	assert.NilError(t, err)
	assert.Equal(t, read.Author.Name, "Renamed", "Refresh returns current data")
}

func TestRepoResolveWithoutFields(t *testing.T) {
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

	err = client.SpecialRelationRepo().Resolve(ctx, read)
	assert.NilError(t, err)
	assert.Assert(t, read.Author.IsPartial(), "without relations nothing is resolved")

	assert.Error(t,
		client.SpecialRelationRepo().Resolve(ctx, nil, with.SpecialRelation.Author()),
		"the passed node must not be nil",
	)
	assert.Error(t,
		client.SpecialRelationRepo().Resolve(ctx, &model.SpecialRelation{}, with.SpecialRelation.Author()),
		"cannot resolve SpecialRelation without existing record ID",
	)
}
