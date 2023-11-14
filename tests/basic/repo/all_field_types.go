package repo

import (
	"context"
	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/by"
	"github.com/go-surreal/som/tests/basic/gen/som/where"
	"github.com/go-surreal/som/tests/basic/gen/som/with"
	"github.com/go-surreal/som/tests/basic/model"
	"time"
)

type AllFieldTypesRepo interface {
	som.AllFieldTypesRepo

	FindByID(ctx context.Context, id string) (*model.AllFieldTypes, error)
	List(ctx context.Context) ([]*model.AllFieldTypes, error)
}

type user struct {
	som.AllFieldTypesRepo
}

func User(db som.Client) AllFieldTypesRepo {
	return &user{
		AllFieldTypesRepo: db.AllFieldTypesRepo(),
	}
}

func (r *user) FindByID(ctx context.Context, id string) (*model.AllFieldTypes, error) {
	return r.AllFieldTypesRepo.Query().Filter(where.AllFieldTypes.ID.Equal(id)).First(ctx)
}

func (r *user) FetchByID(
	ctx context.Context,
	id string,
	fetch ...with.Fetch_[model.AllFieldTypes],
) (
	[]*model.AllFieldTypes,
	error,
) {
	return r.AllFieldTypesRepo.Query().
		Filter(where.AllFieldTypes.ID.Equal(id)).
		Fetch(fetch...).
		All(ctx)
}

var b byte

func (r *user) List(ctx context.Context) ([]*model.AllFieldTypes, error) {
	return r.AllFieldTypesRepo.Query().
		Filter(
			where.Any[model.AllFieldTypes](
				// where.User.ID.Equal("9rb97n04ggwmekxats5a"),
				// where.User.ID.Equal("lvsl8w9gx5i97vado4tp"),
				// where.User.MainGroup().ID.Equal("wq4p7fj4efocis35znzz"),
				// where.User.MyGroups().Since.Before(time.Now()), // ->(member_of where since < $)
				// where.User.MyGroups().Group().ID.Equal(""),     // ->member_of->(group where id = $)

				where.AllFieldTypes.MemberOf().Group().Members().User(
					where.AllFieldTypes.ID.Equal("klkl4w6i9z8u0uyo5w7f"),
				),

				where.AllFieldTypes.Byte.Equal(b),

				where.AllFieldTypes.BytePtr.Equal(b),

				where.AllFieldTypes.ByteSlice().Equal([]byte("omr4f")),

				where.AllFieldTypes.ByteSlicePtr().Equal([]byte("")),

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
			by.AllFieldTypes.CreatedAt.Asc(),
			by.AllFieldTypes.MainGroup().Name.Asc(),
		).
		Fetch(
			with.AllFieldTypes, // this is implicit
			with.AllFieldTypes.MainGroup(),
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
