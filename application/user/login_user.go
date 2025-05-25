package user

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"test_auth/application/ports"
)

type LoginUser struct {
	Repo         ports.UserRepository
	TokenService ports.JwtProvider
}

// doc: responsible for finding the user in the repository
// relied on Repo.FindByUsername method make Execute method works properly
//
// LoginUser is inheriting from UserRepository interface AND JWT provider interface
// Logic Flow Register Interface -> Class inheriting from UserRepository -> Repo.FindByName method -> Execute method

func (uc *LoginUser) Execute(username, password string) (string, error) {
	user, err := uc.Repo.FindByUsername(username)
	if err != nil {
		return "", fmt.Errorf("user %s not found", username)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}
	return uc.TokenService.Generate(user.ID)
}
