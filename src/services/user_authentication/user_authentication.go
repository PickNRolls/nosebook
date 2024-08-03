package services

import (
	"errors"
	"nosebook/src/domain/users"
	"nosebook/src/services/user_authentication/commands"
	"nosebook/src/services/user_authentication/interfaces"

	"golang.org/x/crypto/bcrypt"
)

type UserAuthenticationService struct {
	repo interfaces.UserRepository
}

func NewUserAuthenticationService(repo interfaces.UserRepository) *UserAuthenticationService {
	return &UserAuthenticationService{
		repo: repo,
	}
}

func (s *UserAuthenticationService) RegisterUser(c *commands.RegisterUserCommand) (*users.User, error) {
	existingUser := s.repo.FindByNick(c.Nick)
	if existingUser != nil {
		return nil, errors.New("Can't register user with such nickname.")
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	user := users.NewUser(c.FirstName, c.LastName, c.Nick, string(passhash))
	user, err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserAuthenticationService) LoginUser() (*users.User, error) {
	return users.NewUser("", "", "", ""), nil
}

func (s *UserAuthenticationService) LogoutUser() (*users.User, error) {
	return users.NewUser("", "", "", ""), nil
}
