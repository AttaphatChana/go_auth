package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"test_auth/application/ports"
	"test_auth/domain"
)

type RegisterUser struct {
	Repo ports.UserRepository
}

// doc: responsible for generating a hashed password and saving a new user to the repository
// relied on Repo.Save method to persist the user data to the DB
// need to implement Repo.Save method in the UserRepository interface in order for Execute to work
// Logic Flow Register Interface -> Class inheriting from UserRepository -> Repo.Save method -> Execute method

func (uc *RegisterUser) Execute(username, password string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		ID:       uuid.NewString(),
		Username: username,
		Password: string(hash),
	}
	if err := uc.Repo.Save(user); err != nil {
		return nil, err
	}
	return user, nil
}
