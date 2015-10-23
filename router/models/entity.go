package models

import "github.com/FoxComm/libs/spree"

type Entity struct {
	Id   int
	Name string
	Type string
}

func NewCauseEntity(cause spree.Cause) Entity {
	entity := Entity{
		Id:   cause.Id,
		Name: cause.Name,
		Type: "cause",
	}
	return entity
}

func NewUserEntity(user spree.User) Entity {
	entity := Entity{
		Id:   user.Id,
		Name: user.Name(),
		Type: "user",
	}
	return entity
}

func (e Entity) Parents() []Entity {
	return []Entity{}
}

func (e Entity) Children() []Entity {
	return []Entity{}
}
