package relate

import (
	"errors"
	conv "github.com/marcbinz/som/example/gen/som/conv"
	model "github.com/marcbinz/som/example/model"
	surrealdbgo "github.com/surrealdb/surrealdb.go"
)

type memberOf struct {
	db Database
}

func (e memberOf) Create(edge *model.MemberOf) error {
	if edge.ID != "" {
		return errors.New("ID must not be set for an edge to be created")
	}
	if edge.User.ID == "" {
		return errors.New("ID of the incoming node 'User' must not be empty")
	}
	if edge.Group.ID == "" {
		return errors.New("ID of the outgoing node 'Group' must not be empty")
	}
	query := "RELATE "
	query += "user:" + edge.User.ID
	query += "->member_of->"
	query += "group:" + edge.Group.ID
	query += " CONTENT $data"
	data := conv.FromMemberOf(edge)
	raw, err := e.db.Query(query, map[string]any{"data": data})
	if err != nil {
		return err
	}
	var convEdge conv.MemberOf
	ok, err := surrealdbgo.UnmarshalRaw(raw, &convEdge)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("result is empty")
	}
	*edge = *conv.ToMemberOf(&convEdge)
	return nil
}
func (memberOf) Update(edge *model.MemberOf) error {
	return nil
}
func (memberOf) Delete(edge *model.MemberOf) error {
	return nil
}
