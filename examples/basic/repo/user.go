package repo

import (
	"context"
	"github.com/marcbinz/som/examples/basic/gen/som"
	"github.com/marcbinz/som/examples/basic/gen/som/by"
	"github.com/marcbinz/som/examples/basic/gen/som/where"
	"github.com/marcbinz/som/examples/basic/gen/som/with"
	"github.com/marcbinz/som/examples/basic/model"
	"time"
)

type UserRepo interface {
	som.UserRepo

	FindByID(ctx context.Context, id string) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
}

type user struct {
	som.UserRepo
}

func User(db som.Client) UserRepo {
	return &user{
		UserRepo: db.UserRepo(),
	}
}

func (r *user) FindByID(ctx context.Context, id string) (*model.User, error) {
	return r.UserRepo.Query().Filter(where.User.ID.Equal(id)).First(ctx)
}

func (r *user) FetchByID(
	ctx context.Context,
	id string,
	fetch ...with.Fetch_[model.User],
) (
	[]*model.User,
	error,
) {
	return r.UserRepo.Query().
		Filter(where.User.ID.Equal(id)).
		Fetch(fetch...).
		All(ctx)
}

func (r *user) List(ctx context.Context) ([]*model.User, error) {
	return r.UserRepo.Query().
		Filter(
			where.Any[model.User](
				// where.User.ID.Equal("9rb97n04ggwmekxats5a"),
				// where.User.ID.Equal("lvsl8w9gx5i97vado4tp"),
				// where.User.MainGroup().ID.Equal("wq4p7fj4efocis35znzz"),
				// where.User.MyGroups().Since.Before(time.Now()), // ->(member_of where since < $)
				// where.User.MyGroups().Group().ID.Equal(""),     // ->member_of->(group where id = $)

				where.User.MemberOf().Group().Members().User(
					where.User.ID.Equal("klkl4w6i9z8u0uyo5w7f"),
				),

				//
				// where.User.Groups().ID.In(nil),
				// where.User.Groups().Name.In(nil),

				// where.User.Groups(
				// 	where.Group.ID.Equal(""),
				// 	where.Group.Name.Equal(""),
				// ),

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
		Fetch(
			with.User, // this is implicit
			with.User.MainGroup(),
		).
		// FetchRecordsWithDepth(3).
		// FetchEdgesWithDepth3().
		// FetchDepth(2).
		// OrderRandom().
		// Offset(10).
		// Limit(10).
		Timeout(10 * time.Second).
		Parallel(true).
		// Select( /* choose what data (apart from basic node) should be loaded into the model */ ).
		All(ctx)
}
