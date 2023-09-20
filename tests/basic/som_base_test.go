package basic

import (
	"context"
	sombase "github.com/go-surreal/som"
	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/where"
	"github.com/go-surreal/som/tests/basic/model"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"
	"math"
	"testing"
	"time"
	"unicode/utf8"
)

const (
	surrealDBContainerVersion = "1.0.0"
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
	}
}

func TestQuery(t *testing.T) {
	client := &som.ClientImpl{}

	query := client.AllFieldTypesRepo().Query().
		Filter(
			// where.User.Groups(
			// 	where.Group.Name.Equal(""),
			// 	where.Group.CreatedAt.In(nil),
			// ).Count().GreaterThan(3),
			// where.User.Groups(where.Group.CreatedAt.After(time.Now())),

			where.AllFieldTypes.
				MemberOf(
					where.GroupMember.CreatedAt.Before(time.Now()),
				).
				Group(
					where.Group.ID.Equal("some_id"),
				),

			// select * from user where ->(member_of where createdAt before time::now)->(group where ->(member_of)->(user where id = ""))
			// where.User.MyGroups(where.MemberOf.CreatedAt.Before(time.Now)).Group().Members().User().ID.Equal(""),
		)

	assert.Equal(t,
		"SELECT * FROM all_field_types WHERE (->group_member[WHERE (created_at < $0)]->group[WHERE (id = $1)])",
		// "SELECT * FROM user WHERE (count(groups[WHERE (name = $0 AND created_at INSIDE $1)]) > $2 "+
		// 	"AND groups[WHERE (created_at > $3)]) ",
		query.Describe(),
	)

	//  ("SELECT * FROM user WHERE count(groups[WHERE name = $0]) > $1 " + "AND groups[WHERE created_at > $2].created_at < $3" string)

	// query = query.Filter(
	// 	where.Any(
	// 		where.User.TimePtr.Nil(),
	// 		where.User.UUID.Equal(uuid.New()),
	// 	),
	// )
	//
	// assert.Equal(t,
	// 	"SELECT * FROM user WHERE (string INSIDE $0 AND bool == $1 AND (time_ptr == $2 OR uuid = $3)) ",
	// 	query.Describe(),
	// )
}

func TestWithDatabase(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	str := "Some User"
	uid := uuid.New()

	userNew := model.AllFieldTypes{
		String: str,
		UUID:   uid,
	}

	userIn := userNew

	err := client.AllFieldTypesRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err := client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.ID.Equal(userIn.ID()),
			where.AllFieldTypes.String.Equal(str),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, str, userOut.String)
	assert.Equal(t, uid, userOut.UUID)

	assert.DeepEqual(t,
		userNew, *userOut,
		cmpopts.IgnoreUnexported(sombase.Node{}, sombase.Timestamps{}),
	)
}

func TestNumerics(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	str := "user"

	// MAX

	userMax := model.AllFieldTypes{
		String: str,
		Int:    math.MaxInt,
		Int8:   math.MaxInt8,
		Int16:  math.MaxInt16,
		Int32:  math.MaxInt32,
		Int64:  math.MaxInt64,
		// Uint:    math.MaxUint,
		Uint8:  math.MaxUint8,
		Uint16: math.MaxUint16,
		Uint32: math.MaxUint32,
		// Uint64:  0, // math.MaxUint64,
		// Uintptr: 0, // math.MaxUint64,
		Float32: math.MaxFloat32,
		Float64: math.MaxFloat64,
		Rune:    math.MaxInt32,
	}

	userIn := userMax

	err := client.AllFieldTypesRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err := client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.ID.Equal(userIn.ID()),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t,
		userMax, *userOut,
		cmpopts.IgnoreUnexported(sombase.Node{}, sombase.Timestamps{}),
	)

	// MIN

	userMin := model.AllFieldTypes{
		String: str,
		Int:    math.MinInt,
		Int8:   math.MinInt8,
		Int16:  math.MinInt16,
		Int32:  math.MinInt32,
		Int64:  math.MinInt64,
		// Uint:    math.MaxUint,
		Uint8:  0,
		Uint16: 0,
		Uint32: 0,
		// Uint64:  0, // math.MaxUint64,
		// Uintptr: 0, // math.MaxUint64,
		Float32: -math.MaxFloat32,
		Float64: -math.MaxFloat64,
		Rune:    math.MinInt32,
	}

	userIn = userMin

	err = client.AllFieldTypesRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.ID.Equal(userIn.ID()),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t,
		userMin, *userOut,
		cmpopts.IgnoreUnexported(sombase.Node{}, sombase.Timestamps{}),
	)
}

func TestTimestamps(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	user := &model.AllFieldTypes{}

	err := client.AllFieldTypesRepo().Create(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, !user.CreatedAt().IsZero())
	assert.Check(t, !user.UpdatedAt().IsZero())
	assert.Check(t, time.Since(user.CreatedAt()) < time.Second)
	assert.Check(t, time.Since(user.UpdatedAt()) < time.Second)

	time.Sleep(time.Second)

	err = client.AllFieldTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, !user.CreatedAt().IsZero())
	assert.Check(t, !user.UpdatedAt().IsZero())
	assert.Check(t, time.Since(user.CreatedAt()) > time.Second)
	assert.Check(t, time.Since(user.UpdatedAt()) < time.Second)
}

func FuzzWithDatabase(f *testing.F) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, f)
	defer cleanup()

	f.Add("Some User")

	f.Fuzz(func(t *testing.T, str string) {
		userIn := &model.AllFieldTypes{
			String: str,
		}

		err := client.AllFieldTypesRepo().Create(ctx, userIn)
		if err != nil {
			t.Fatal(err)
		}

		if userIn.ID() == "" {
			t.Fatal("user ID must not be empty after create call")
		}

		userOut, err := client.AllFieldTypesRepo().Query().
			Filter(
				where.AllFieldTypes.ID.Equal(userIn.ID()),
			).
			First(ctx)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, userIn.String, userOut.String)
	})
}

func FuzzCustomModelIDs(f *testing.F) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, f)
	defer cleanup()

	f.Add("v9uitj942tv2403tnv")
	f.Add("vb92thj29v4tjn20d3")
	f.Add("ij024itvnjc20394it")
	f.Add(" 0")
	f.Add("\"0")
	f.Add("🙂")
	f.Add("✅")
	f.Add("👋😉")

	f.Fuzz(func(t *testing.T, id string) {
		if !utf8.ValidString(id) {
			t.Skip("id is not a valid utf8 string")
		}

		userIn := &model.AllFieldTypes{
			String: "1",
		}

		err := client.AllFieldTypesRepo().CreateWithID(ctx, id, userIn)
		if err != nil {
			t.Fatal(err)
		}

		if userIn.ID() == "" {
			t.Fatal("user ID must not be empty after create call")
		}

		userOut, ok, err := client.AllFieldTypesRepo().Read(ctx, userIn.ID())

		if err != nil {
			t.Fatal(err)
		}

		if !ok {
			t.Fatalf("user with ID '%s' not found", userIn.ID())
		}

		assert.Equal(t, userIn.ID(), userOut.ID())
		assert.Equal(t, "1", userOut.String)

		userOut.String = "2"

		err = client.AllFieldTypesRepo().Update(ctx, userOut)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "2", userOut.String)

		err = client.AllFieldTypesRepo().Delete(ctx, userOut)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func BenchmarkWithDatabase(b *testing.B) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, b)
	defer cleanup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		userIn := &model.AllFieldTypes{
			String: "Some User",
		}

		err := client.AllFieldTypesRepo().Create(ctx, userIn)
		if err != nil {
			b.Fatal(err)
		}

		if userIn.ID() == "" {
			b.Fatal("user ID must not be empty after create call")
		}

		userOut, err := client.AllFieldTypesRepo().Query().
			Filter(
				where.AllFieldTypes.ID.Equal(userIn.ID()),
			).
			First(ctx)

		if err != nil {
			b.Fatal(err)
		}

		assert.Equal(b, userIn.String, userOut.String)
	}
}

func TestAsync(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{})
	if err != nil {
		t.Fatal(err)
	}

	resCh := client.AllFieldTypesRepo().Query().
		Filter().
		CountAsync(ctx)

	assert.NilError(t, <-resCh.Err())
	assert.Equal(t, 1, <-resCh.Val())

	err = client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{})
	if err != nil {
		t.Fatal(err)
	}

	resCh = client.AllFieldTypesRepo().Query().
		Filter().
		CountAsync(ctx)

	assert.NilError(t, <-resCh.Err())
	assert.Equal(t, 2, <-resCh.Val())
}
