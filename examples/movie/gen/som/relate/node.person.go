// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package relate

func NewPerson(db Database) *Person {
	return &Person{db: db}
}

type Person struct {
	db Database
}
