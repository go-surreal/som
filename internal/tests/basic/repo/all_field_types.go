package repo

import (
	"context"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/by"
	"github.com/go-surreal/som/tests/basic/gen/som/repo"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/gen/som/with"
	"github.com/go-surreal/som/tests/basic/model"
)

type AllFieldTypesRepo interface {
	repo.AllFieldTypesRepo

	FindByID(ctx context.Context, id *som.ID) (*model.AllFieldTypes, error)
	List(ctx context.Context) ([]*model.AllFieldTypes, error)
}

type allFieldTypesRepo struct {
	repo.AllFieldTypesRepo
}

func NewAllFieldTypesRepo(db repo.Client) AllFieldTypesRepo {
	return &allFieldTypesRepo{
		AllFieldTypesRepo: db.AllFieldTypesRepo(),
	}
}

func (r *allFieldTypesRepo) FindByID(ctx context.Context, id *som.ID) (*model.AllFieldTypes, error) {
	return r.AllFieldTypesRepo.Query().Where(filter.AllFieldTypes.ID.Equal(id)).First(ctx)
}

func (r *allFieldTypesRepo) FetchByID(
	ctx context.Context,
	id *som.ID,
	fetch ...with.Fetch_[model.AllFieldTypes],
) (
	[]*model.AllFieldTypes,
	error,
) {
	return r.AllFieldTypesRepo.Query().
		Where(filter.AllFieldTypes.ID.Equal(id)).
		Fetch(fetch...).
		All(ctx)
}

var b byte

func (r *allFieldTypesRepo) List(ctx context.Context) ([]*model.AllFieldTypes, error) {
	return r.AllFieldTypesRepo.Query().
		Where(
			//filter.Any[model.AllFieldTypes](
			//	filter.User.ID.Equal("9rb97n04ggwmekxats5a"),
			//	filter.User.ID.Equal("lvsl8w9gx5i97vado4tp"),
			//	filter.User.MainGroup().ID.Equal("wq4p7fj4efocis35znzz"),
			//	filter.User.MyGroups().Since.Before(time.Now()), // ->(member_of where since < $)
			//	filter.User.MyGroups().Group().ID.Equal(""),     // ->member_of->(group where id = $)
			//
			//	filter.AllFieldTypes.MemberOf().Group().Members().User(
			//		filter.AllFieldTypes.ID.Equal("klkl4w6i9z8u0uyo5w7f"),
			//	),
			//
			//	filter.AllFieldTypes.Byte.Equal(b),
			//
			//	filter.AllFieldTypes.BytePtr.Equal(b),
			//
			//	filter.AllFieldTypes.ByteSlice().Equal([]byte("omr4f")),
			//
			//	filter.AllFieldTypes.ByteSlicePtr().Equal([]byte("")),
			//
			//	filter.User.Groups().ID.In(nil),
			//	filter.User.Groups().Name.In(nil),
			//
			//	filter.User.Groups(
			//		filter.Group.ID.Equal(""),
			//		filter.Group.Name.Equal(""),
			//	),
			//),
			// filter.User.String.FuzzyMatch("my fuzzy value"),
			// filter.User.UUID.Equal(uuid.UUID{}),
			// filter.Any(
			// 	filter.User.Role.Equal(""),
			// 	filter.User.CreatedAt.Before(time.Now()),
			// ),
			// filter.All(
			// 	filter.User.Role.Equal(""),
			// 	filter.User.Groups().Name.FuzzyMatch("some group"),
			// 	filter.User.Groups().Contains(model.Group{}),
			// ),
			// filter.User.ID.Equal(""),
			// filter.User.Login().Username.Equal(""),
			// filter.User.Role.Equal(""),
			// filter.User.Groups().Count().GreaterThan(5),
			// //
			// filter.User.Other().Contains(""),
			// filter.User.Other().ContainsAll([]string{"", ""}),
			// filter.User.Roles().ContainsNot(model.RoleAdmin),
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
