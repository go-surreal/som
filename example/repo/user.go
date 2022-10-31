package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/marcbinz/sdb/example/gen/sdb"
	"github.com/marcbinz/sdb/example/gen/sdb/by"
	"github.com/marcbinz/sdb/example/gen/sdb/where"
	"github.com/marcbinz/sdb/example/model"
	"time"
)

type User struct {
	DB *sdb.Client
}

func (repo *User) Create(ctx context.Context, user *model.User) error {
	return sdb.User.Create(ctx, repo.DB, user)
}

func (repo *User) FindById(ctx context.Context, id string) (*model.User, error) {
	return sdb.User.Read(ctx, repo.DB, id)
}

func (repo *User) List(ctx context.Context) ([]*model.User, error) {
	return sdb.User.Query().
		Filter(
			where.User.String.FuzzyMatch("my fuzzy value"),
			where.User.UUID.Equal(uuid.UUID{}),
			where.Any(
				where.User.Role.Equal(""),
				where.User.CreatedAt.Before(time.Now()),
			),
			where.All(
				where.User.Role.Equal(""),
				where.User.Groups().Name.FuzzyMatch("some group"),
				where.User.Groups().Contains(model.Group{}),
			),
			where.User.ID.Equal(""),
			where.User.Login().Username.Equal(""),
			where.User.Role.Equal(""),
			where.User.Groups().Count().GreaterThan(5),
			//
			where.User.Other().Contains(""),
			where.User.Other().ContainsAll([]string{"", ""}),
			where.User.Roles().ContainsNot(model.RoleAdmin),
		).
		Sort(
			by.User.ID.Asc(),
			by.User.String.Collate().Desc(),
			by.User.String.Numeric().Asc(),
			by.User.CreatedAt.Asc(),
			by.User.CreatedAt.Desc(),
			by.User.Random(),
		).
		Offset(10).
		Limit(10).
		Timeout(10 * time.Second).
		Parallel(true).
		// Select( /* choose what data (apart from basic node) should be loaded into the model */ ).
		All()
}
