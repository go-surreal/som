package fulltext

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/google/go-cmp/cmp/cmpopts"
	sombase "github.com/marcbinz/som"
	"github.com/marcbinz/som/examples/basic/gen/som"
	"github.com/marcbinz/som/examples/basic/gen/som/where"
	"github.com/marcbinz/som/examples/basic/model"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/v3/assert"
	"testing"
)

const (
	surrealDBContainerVersion = "nightly"
	containerName             = "som_test_surrealdb"
	containerStartedMsg       = "Started web server on 0.0.0.0:8000"
)

func conf(endpoint string) som.Config {
	return som.Config{
		Address:   "ws://" + endpoint,
		Username:  "root",
		Password:  "root",
		Namespace: "som_test",
		Database:  "example_basic",
		Init:      dbInit,
	}
}

func dbInit(db som.InitDB) {

	// DEFINE ANALYZER simple TOKENIZERS camel,class FILTERS snowball(English)
	analyzerSimple := db.NewAnalyzer("simple").
		Tokenizers(som.Camel, som.Class).  // som.Blank, som.Camel, som.Class, som.Punct
		Filters(som.Snowball(som.English)) // som.Ascii, som.EdgeNGram, som.Snowball(Lang), som.Lowercase, som.Uppercase

	// DEFINE INDEX idx_author ON books FIELDS author
	db.Table().Books().NewIndex("idx_author").
		Fields(som.TableBooks.Fields.Author)

	// DEFINE INDEX uniq_isbn ON books FIELDS isbn UNIQUE
	db.Table().Books().NewIndex("uniq_isbn").
		Fields(som.TableBooks.Fields.ISBN).
		Unique()

	// DEFINE INDEX ft_title ON books FIELDS title SEARCH analyzer_simple BM25(1.2, 0.75) HIGHLIGHTS
	db.Table().Books().NewIndex("ft_title").
		Fields(som.TableBooks.Fields.Title).
		Search(analyzerSimple, som.BM25(1.2, 0.75), som.Highlights)

}

func TestSearch(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	str := "Some User"

	userNew := model.User{
		String: str,
	}

	userIn := userNew

	err := client.UserRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err := client.UserRepo().Query().
		Filter(
			where.User.ID.Equal(userIn.ID()),
			where.User.String.Equal(str),
		).
		Fulltext(
			search.User.String("me us"),
		).
		WithHighlights().
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, str, userOut.String)

	assert.DeepEqual(t,
		userNew, *userOut,
		cmpopts.IgnoreUnexported(sombase.Node{}, sombase.Timestamps{}),
	)
}

//
// -- HELPER
//

func prepareDatabase(ctx context.Context, tb testing.TB) (som.Client, func()) {
	tb.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	req := testcontainers.ContainerRequest{
		Name:         containerName,
		Image:        "surrealdb/surrealdb:" + surrealDBContainerVersion,
		Cmd:          []string{"start", "--strict", "--allow-funcs", "--user", "root", "--pass", "root", "--log", "debug", "memory"},
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForLog(containerStartedMsg),
		HostConfigModifier: func(conf *container.HostConfig) {
			conf.AutoRemove = true
		},
	}

	surreal, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Reuse:            true,
		},
	)
	if err != nil {
		tb.Fatal(err)
	}

	endpoint, err := surreal.Endpoint(ctx, "")
	if err != nil {
		tb.Fatal(err)
	}

	client, err := som.NewClient(ctx, conf(endpoint))
	if err != nil {
		tb.Fatal(err)
	}

	if err := client.ApplySchema(ctx); err != nil {
		tb.Fatal(err)
	}

	cleanup := func() {
		client.Close()

		if err := surreal.Terminate(ctx); err != nil {
			tb.Fatalf("failed to terminate container: %s", err.Error())
		}
	}

	return client, cleanup
}
