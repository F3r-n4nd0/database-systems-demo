package user

import "web-service-users/src/model"

type WebServiceKudos interface {
	GetKudosByUserName(userName string) ([]*model.Kudos, error)
}
