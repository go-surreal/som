package example

import (
	"github.com/marcbinz/sdb/define"
	"github.com/marcbinz/sdb/define/field"
)

func Schema() *define.Schema {

	def := define.New()

	login := def.Object("login").
		With(
			field.String("username").Pointer(),
			field.String("password").Pointer(),
		)

	user := def.Table("user").
		With(
			field.String("name").Pointer(),
			field.Object("login", login),
			field.Link("main_group", "group").Pointer(),
			field.Link("groups", "group").Pointer().Slice(),
		)

	post := def.Table("post").
		With(
			field.String("title").Pointer(),
		)

	user.With(field.Link("posts", post))

	def.Edge("wrote").From(user).To(post).
		With(
			field.Time("written_at"),
		)

	return def.Schema()
}
