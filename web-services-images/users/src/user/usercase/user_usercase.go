package usercase

import (
	"time"
	"web-service-users/src/model"
	"web-service-users/src/user"
)

type userUserCase struct {
	userRepos user.Repository
}

func NewUserUserCase(u user.Repository) user.UseCase {
	return &userUserCase{
		userRepos: u,
	}
}

func (u *userUserCase) Store(user *model.User) (*model.User, error) {

	existedUser, _ := u.userRepos.GetByUserName(user.UserName)
	if existedUser != nil {
		return nil, model.ConflictError
	}
	user.QuantityKudos = 0
	user.CreateAt = time.Now()
	user.UpdateAt = time.Now()
	err := u.userRepos.Store(user)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (u *userUserCase) GetByUserName(userName string) (*model.User, error) {

	res, err := u.userRepos.GetByUserName(userName)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (u *userUserCase) DeleteByUserName(userName string) error {

	err := u.userRepos.DeleteByUserName(userName)
	if err != nil {
		return err
	}
	return nil

}

func (u *userUserCase) FetchAllUsers(pageSize int64, numberPage int64) ([]*model.User, error) {

	users, err := u.userRepos.FetchAllUsers(pageSize, numberPage)
	if err != nil {
		return nil, err
	}
	return users, nil

}

func (u *userUserCase) UpdateQuantityKudos(userName string, quantity int) error {

	err := u.userRepos.UpdateQuantityKudos(userName, int32(quantity), time.Now())
	if err != nil {
		return err
	}
	return nil

}
