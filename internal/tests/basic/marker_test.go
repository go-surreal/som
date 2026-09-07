package basic

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
	"som.test/gen/som"
	"som.test/gen/som/with"
	"som.test/model"
)

func TestMarkerLoaded(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := model.SpecialTypes{Name: "Author"}

	assert.Equal(t, record.Marker(), som.Marker(0), "a model built by the application has no marker")
	assert.Assert(t, !record.IsPartial())

	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)

	read, exists, err := client.SpecialTypesRepo().Read(ctx, string(record.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)

	assert.Assert(t, read.Marker().Has(som.MarkerLoaded))
	assert.Assert(t, !read.IsPartial())
}

func TestMarkerPartialLink(t *testing.T) {
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
	assert.Assert(t, read.Author != nil)

	assert.Equal(t, read.Author.ID(), author.ID(), "an unfetched link keeps its record id")
	assert.Equal(t, read.Author.Name, "", "an unfetched link holds no field values")
	assert.Assert(t, read.Author.IsPartial())
	assert.Assert(t, read.Author.Marker().Has(som.MarkerLoaded|som.MarkerPartial))

	assert.ErrorIs(t, client.SpecialTypesRepo().Update(ctx, read.Author), som.ErrPartialModel)
	assert.ErrorIs(t, client.SpecialTypesRepo().Delete(ctx, read.Author), som.ErrPartialModel)
	assert.ErrorIs(t, client.SpecialTypesRepo().Erase(ctx, read.Author), som.ErrPartialModel)
	assert.ErrorIs(t, client.SpecialTypesRepo().Restore(ctx, read.Author), som.ErrPartialModel)

	err = client.SpecialTypesRepo().Refresh(ctx, read.Author)
	assert.NilError(t, err)

	assert.Assert(t, !read.Author.IsPartial(), "a refreshed link is fully loaded")
	assert.Equal(t, read.Author.Name, author.Name)

	err = client.SpecialTypesRepo().Update(ctx, read.Author)
	assert.NilError(t, err)
}

func TestMarkerFetchedLink(t *testing.T) {
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

	posts, err := client.SpecialRelationRepo().Query().
		Fetch(with.SpecialRelation.Author()).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(posts), 1)
	assert.Assert(t, posts[0].Author != nil)

	assert.Equal(t, posts[0].Author.Name, author.Name)
	assert.Assert(t, !posts[0].Author.IsPartial(), "a fetched link is fully loaded")
	assert.Assert(t, posts[0].Author.Marker().Has(som.MarkerLoaded))
}
