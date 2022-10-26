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
			where.User.Text.Contains(""),
			where.Any(
				where.User.Role.Equal(""),
			),
			where.All(
				where.User.Role.Equal(""),
			),
			where.User.ID.Equal(""),
			// where.User.Login().Username.Equal(""),
			where.User.Role.Equal(""),
			where.Count(where.User.Groups()).GT(5),
		).
		Sort(
			by.User.ID.Asc(),
			by.User.Role.Collate().Desc(),
			by.User.CreatedAt.Asc(),
			by.User.CreatedAt.Desc(),
			by.User.CreatedAt.Collate().Asc(),  // only strings!
			by.User.CreatedAt.Collate().Desc(), // only strings!
			by.User.CreatedAt.Numeric().Asc(),  // only strings!
			by.User.CreatedAt.Numeric().Desc(), // only strings!
			// by.Rand[predicate.User](), ?!
		).
		Offset(10).
		Limit(10).
		Timeout(10 * time.Second).
		Parallel(true).
		// Select( /* choose what data (apart from basic node) should be loaded into the model */ ).
		All()
}
