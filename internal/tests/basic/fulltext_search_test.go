package basic

import (
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som/where"
	"github.com/go-surreal/som/tests/basic/model"
)

func TestFullTextSearch(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "this is a test string",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	query := client.AllFieldTypesRepo().Query().
		Search(
			where.AllFieldTypes.String.Matches("this is a test"),
			where.AllFieldTypes.String.Equal(""),
		).
		Filter(
			where.AllFieldTypes.String.Matches("this is a test"),
			where.AllFieldTypes.String.Equal(""),
			where.AllFieldTypes.Login().Password.Matches(""),
		)

	t.Logf("%s", query.Describe())

	res, err := query.AllMatches(ctx)
	if err != nil {
		t.Fatalf("failed to execute full-text search query: %v", err)
	}

	for _, item := range res {
		t.Logf("Found item: %+v", item)
	}
}
