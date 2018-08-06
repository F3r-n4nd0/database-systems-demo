package usercase

import (
	"stats/src/kudos"
	"stats/src/model"

	log "github.com/sirupsen/logrus"
)

type kudosUserCase struct {
	KudosWebService kudos.WebServiceKudos
	UsersWebService kudos.WebServiceUsers
}

func NewKudosUserCase(uws kudos.WebServiceUsers, kws kudos.WebServiceKudos) kudos.UseCase {
	return &kudosUserCase{kws, uws}
}

func (k *kudosUserCase) QueueKudos(kudosQ *model.KudosQueue) {

	currentQuantity, err := k.KudosWebService.GetQuantityKudos(kudosQ.UserName)
	if err != nil {
		log.Error(err)
		return
	}
	err = k.UsersWebService.UpdateQuantityKudos(kudosQ.UserName, currentQuantity)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Kudos quantity updated")

}
