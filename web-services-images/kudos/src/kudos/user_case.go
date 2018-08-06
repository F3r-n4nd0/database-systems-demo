package kudos

import "web-service-kudos/src/model"

type UseCase interface {
	Store(kudos *model.Kudos) (*model.Kudos, error)
	GetByID(id string) (*model.Kudos, error)
	DeleteByID(id string) error
	FetchAllKudos(pageSize int64, numberPage int64) ([]*model.Kudos, error)
	GetQuantityByUserName(userName string) (int, error)
}
