package user

import (
	"time"
	"web-service-users/src/model"
)

type Repository interface {
	Store(user *model.User) error
	GetByUserName(userName string) (*model.User, error)
	DeleteByUserName(userName string) error
	FetchAllUsers(pageSize int64, numberPage int64) ([]*model.User, error)
	UpdateQuantityKudos(userName string, quantity int32, updateDate time.Time) error
}
