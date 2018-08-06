package usercase

import (
	"time"
	"web-service-kudos/src/kudos"
	"web-service-kudos/src/model"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type kudosUserCase struct {
	kudosRepos kudos.Repository
	KudosQueue kudos.Queue
}

func NewKudosUsecase(k kudos.Repository, kq kudos.Queue) kudos.UseCase {
	return &kudosUserCase{
		kudosRepos: k,
		KudosQueue: kq,
	}
}

func (k *kudosUserCase) Store(kudos *model.Kudos) (*model.Kudos, error) {

	existedKudos, _ := k.kudosRepos.GetByID(kudos.Id)
	if existedKudos != nil {
		return nil, model.ConflictError
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	kudos.Id = uuid.String()
	kudos.CreateAt = time.Now()
	kudos.UpdateAt = time.Now()
	err = k.kudosRepos.Store(kudos)
	if err != nil {
		return nil, err
	}
	err = k.KudosQueue.IncreasesKudos(kudos)
	if err != nil {
		log.Error(err)
	}
	return kudos, nil

}

func (k *kudosUserCase) GetByID(id string) (*model.Kudos, error) {

	res, err := k.kudosRepos.GetByID(id)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (k *kudosUserCase) DeleteByID(id string) error {

	kudos, err := k.GetByID(id)
	if err != nil {
		return err
	}
	err = k.kudosRepos.DeleteByID(id)
	if err != nil {
		return err
	}
	err = k.KudosQueue.DecreasesKudos(kudos)
	if err != nil {
		log.Error(err)
	}
	return nil

}

func (k *kudosUserCase) FetchAllKudos(pageSize int64, numberPage int64) ([]*model.Kudos, error) {

	kudos, err := k.kudosRepos.FetchAllKudos(pageSize, numberPage)
	if err != nil {
		return nil, err
	}
	return kudos, nil

}

func (k *kudosUserCase) GetQuantityByUserName(userName string) (int, error) {

	quantityKudos, err := k.kudosRepos.GetQuantityByUserName(userName)
	if err != nil {
		return 0, err
	}
	return quantityKudos, nil

}
