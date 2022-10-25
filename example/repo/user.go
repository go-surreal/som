package repo

import (
	"context"
	"github.com/marcbinz/sdb/example/gen/sdb"
	"github.com/marcbinz/sdb/example/gen/sdb/where"
	"github.com/marcbinz/sdb/example/model"
)

type User struct {
	db sdb.Client
}

func (repo *User) List(ctx context.Context) ([]*model.User, error) {

	return sdb.User.Query().
		Filter(
			where.User.ID.Equal(""),
			where.User.Username.Equal(""),
			where.User.Role.Equal(model.Role{}),
		).
		All()

}
