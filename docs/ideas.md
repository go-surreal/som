
```go
package som

func Select[M any](ctx context.Context, handler Query[M], field with.Fetch_[M]) {
	
}

type Query[M any] interface {
	All([]*M, error)
}

type Repository[M any]interface{
	BaseRepository[M]
	
	Query() Query[M]
}

type BaseRepository[M any]interface{
	Create(ctx context.Context, model *M) error
	Update(ctx context.Context, model *M) error
	Delete(ctx context.Context, model *M) error
}

```


```go

repo.db.User.Modify().
	Filter().
	Set(
	  field.User.Name.Value("new name"),
    ).
	Execute(ctx)

```
