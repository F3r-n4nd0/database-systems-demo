package user

import (
	"web-service-users/src/model"
)

type UseCase interface {
	Store(user *model.User) (*model.User, error)
	GetByUserName(userName string) (*model.User, error)
	DeleteByUserName(userName string) error
	FetchAllUsers(pageSize int64, numberPage int64) ([]*model.User, error)
	UpdateQuantityKudos(userName string, quantity int) error
}
