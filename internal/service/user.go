package service

import (
	"github.com/olekturbo/mysterious/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type userStorage interface {
	CreateUser(user *storage.User) error
	FindUser(email string) (*storage.User, error)
}

type idService interface {
	GenerateID() string
}

type User struct {
	userStorage userStorage
	idService   idService
}

func NewUser(userStorage userStorage, idService idService) *User {
	return &User{
		userStorage: userStorage,
		idService:   idService,
	}
}

type CreateParams struct {
	Email    string
	Password string
}

func (u *User) Create(params CreateParams) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.userStorage.CreateUser(&storage.User{
		ID:       u.idService.GenerateID(),
		Email:    params.Email,
		Password: string(hashedPassword),
	})
}

type MatchParams struct {
	Email    string
	Password string
}

type MatchResult struct {
	ID string
}

func (u *User) Match(params MatchParams) (MatchResult, error) {
	user, err := u.userStorage.FindUser(params.Email)
	if err != nil {
		return MatchResult{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		return MatchResult{}, err
	}

	return MatchResult{ID: user.ID}, nil
}
