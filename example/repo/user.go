package repo

import (
	"context"
	"github.com/marcbinz/sdb/example/gen/sdb"
	"github.com/marcbinz/sdb/example/gen/sdb/by"
	"github.com/marcbinz/sdb/example/gen/sdb/where"
	"github.com/marcbinz/sdb/example/model"
	"time"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	FindById(ctx context.Context, id string) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
}

type user struct {
	db *sdb.Client
}

func User(db *sdb.Client) UserRepo {
	return &user{db: db}
}

func (repo *user) Create(ctx context.Context, user *model.User) error {
	return repo.db.User().Create(ctx, user)
}

func (repo *user) FindById(ctx context.Context, id string) (*model.User, error) {
	return repo.db.User().Read(ctx, id)
}

func (repo *user) List(ctx context.Context) ([]*model.User, error) {
	return repo.db.User().Query().
		Filter(
			where.Any(
				where.User.ID.Equal("9rb97n04ggwmekxats5a"),
				where.User.ID.Equal("lvsl8w9gx5i97vado4tp"),
				where.User.MainGroup().ID.Equal("wq4p7fj4efocis35znzz"),
			),
			// where.User.String.FuzzyMatch("my fuzzy value"),
			// where.User.UUID.Equal(uuid.UUID{}),
			// where.Any(
			// 	where.User.Role.Equal(""),
			// 	where.User.CreatedAt.Before(time.Now()),
			// ),
			// where.All(
			// 	where.User.Role.Equal(""),
			// 	where.User.Groups().Name.FuzzyMatch("some group"),
			// 	where.User.Groups().Contains(model.Group{}),
			// ),
			// where.User.ID.Equal(""),
			// where.User.Login().Username.Equal(""),
			// where.User.Role.Equal(""),
			// where.User.Groups().Count().GreaterThan(5),
			// //
			// where.User.Other().Contains(""),
			// where.User.Other().ContainsAll([]string{"", ""}),
			// where.User.Roles().ContainsNot(model.RoleAdmin),
		).
		Order(
			by.User.CreatedAt.Asc(),
			by.User.MainGroup().Name.Asc(),
		).
		// Fetch(
		// 	with.User, // this is implicit
		// 	with.User.MainGroup(),
		// ).
		// FetchRecordsWithDepth(3).
		// FetchEdgesWithDepth3().
		FetchDepth(0).
		// OrderRandom().
		// Offset(10).
		// Limit(10).
		Timeout(10 * time.Second).
		Parallel(true).
		// Select( /* choose what data (apart from basic node) should be loaded into the model */ ).
		All()
}
