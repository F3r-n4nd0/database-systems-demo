package user

import "web-service-users/src/model"

type SearchDB interface {
	FindUser(queryString string) ([]*model.UserFound, error)
}
