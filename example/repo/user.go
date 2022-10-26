package repo

import (
	"context"
	"github.com/marcbinz/sdb/example/gen/sdb"
	"github.com/marcbinz/sdb/example/gen/sdb/by"
	"github.com/marcbinz/sdb/example/gen/sdb/where"
	"github.com/marcbinz/sdb/example/model"
	"time"
)

type User struct {
	db sdb.Client
}

func (repo *User) List(ctx context.Context) ([]*model.User, error) {
	return sdb.User.Query().
		Filter(
			where.User.String.Contains(""),
			where.Any(
				where.User.Role.Equal(""),
				where.User.CreatedAt.Before(time.Now()),
			),
			where.All(
				where.User.Role.Equal(""),
				where.User.MainGroup().Name.Equal("some group"),
			),
			where.User.ID.Equal(""),
			// where.User.Login().Username.Equal(""),
			where.User.Role.Equal(""),
			where.Count(where.User.Groups()).GreaterThan(5),
		).
		Sort(
			by.User.ID.Asc(),
			by.User.String.Collate().Desc(),
			by.User.String.Numeric().Asc(),
			by.User.CreatedAt.Asc(),
			by.User.CreatedAt.Desc(),
			// by.Rand[predicate.User](), ?!
		).
		Offset(10).
		Limit(10).
		Timeout(10 * time.Second).
		Parallel(true).
		// Select( /* choose what data (apart from basic node) should be loaded into the model */ ).
		All()
}
