package with

import model "github.com/marcbinz/sdb/example/model"

var User = user[model.User]("")

type user[T any] string

func (n user[T]) fetch(T) {}
func (n user[T]) MainGroup() group[T] {
	return group[T](keyed(n, "main_group"))
}
