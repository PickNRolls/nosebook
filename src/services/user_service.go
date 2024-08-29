package services

import (
	"errors"
	"nosebook/src/domain/users"
	common_interfaces "nosebook/src/services/common/interfaces"
	"nosebook/src/services/user_service/commands"
)

type UserService struct {
	userRepo common_interfaces.UserRepository
}

func NewUserService(userRepo common_interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (this *UserService) GetUser(c *commands.GetUserCommand) (*users.User, error) {
	user := this.userRepo.FindById(c.Id)
	if user == nil {
		return nil, errors.New("No such user.")
	}

	return user, nil
}

func (this *UserService) GetAllUsers() ([]*users.User, error) {
	all, err := this.userRepo.FindAll()
	if err != nil {
		return make([]*users.User, 0), err
	}

	return all, nil
}
