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

type AllTypesRepo interface {
	repo.AllTypesRepo

	FindByID(ctx context.Context, id *som.ID) (*model.AllTypes, error)
	List(ctx context.Context) ([]*model.AllTypes, error)
}

type allTypesRepo struct {
	repo.AllTypesRepo
}

func NewAllTypesRepo(db repo.Client) AllTypesRepo {
	return &allTypesRepo{
		AllTypesRepo: db.AllTypesRepo(),
	}
}

func (r *allTypesRepo) FindByID(ctx context.Context, id *som.ID) (*model.AllTypes, error) {
	return r.AllTypesRepo.Query().Where(filter.AllTypes.ID.Equal(id)).First(ctx)
}

func (r *allTypesRepo) FetchByID(
	ctx context.Context,
	id *som.ID,
	fetch ...with.Fetch_[model.AllTypes],
) (
	[]*model.AllTypes,
	error,
) {
	return r.AllTypesRepo.Query().
		Where(filter.AllTypes.ID.Equal(id)).
		Fetch(fetch...).
		All(ctx)
}

var b byte

func (r *allTypesRepo) List(ctx context.Context) ([]*model.AllTypes, error) {
	return r.AllTypesRepo.Query().
		Where(
			//filter.Any[model.AllTypes](
			//	filter.User.ID.Equal("9rb97n04ggwmekxats5a"),
			//	filter.User.ID.Equal("lvsl8w9gx5i97vado4tp"),
			//	filter.User.MainGroup().ID.Equal("wq4p7fj4efocis35znzz"),
			//	filter.User.MyGroups().Since.Before(time.Now()), // ->(member_of where since < $)
			//	filter.User.MyGroups().Group().ID.Equal(""),     // ->member_of->(group where id = $)
			//
			//	filter.AllTypes.FieldMemberOf().Group().Members().User(
			//		filter.AllTypes.ID.Equal("klkl4w6i9z8u0uyo5w7f"),
			//	),
			//
			//	filter.AllTypes.FieldByte.Equal(b),
			//
			//	filter.AllTypes.FieldBytePtr.Equal(b),
			//
			//	filter.AllTypes.FieldByteSlice().Equal([]byte("omr4f")),
			//
			//	filter.AllTypes.FieldByteSlicePtr().Equal([]byte("")),
			//
			//	filter.User.Groups().ID.In(nil),
			//	filter.User.Groups().Name.In(nil),
			//
			//	filter.User.Groups(
			//		filter.Group.ID.Equal(""),
			//		filter.Group.Name.Equal(""),
			//	),
			//),
			// filter.User.FieldString.FuzzyMatch("my fuzzy value"),
			// filter.User.FieldUUID.Equal(uuid.UUID{}),
			// filter.Any(
			// 	filter.User.FieldRole.Equal(""),
			// 	filter.User.CreatedAt.Before(time.Now()),
			// ),
			// filter.All(
			// 	filter.User.FieldRole.Equal(""),
			// 	filter.User.Groups().Name.FuzzyMatch("some group"),
			// 	filter.User.Groups().Contains(model.Group{}),
			// ),
			// filter.User.ID.Equal(""),
			// filter.User.FieldCredentials().Username.Equal(""),
			// filter.User.FieldRole.Equal(""),
			// filter.User.Groups().Count().GreaterThan(5),
			// //
			// filter.User.FieldOther().Contains(""),
			// filter.User.FieldOther().ContainsAll([]string{"", ""}),
			// filter.User.FieldRoles().ContainsNot(model.RoleAdmin),
		).
		Order(
			by.AllTypes.CreatedAt.Asc(),
			by.AllTypes.FieldNode().Name.Asc(),
		).
		Fetch(
			with.AllTypes, // this is implicit
			with.AllTypes.FieldNode(),
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
