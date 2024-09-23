package svc

import (
	"untitled/internal/users/mdl"
	"untitled/internal/users/repo"
)

type UserSvc struct {
	userRepo *repo.UserRepo
}

func NewUserSvc(userRepo *repo.UserRepo) *UserSvc {
	return &UserSvc{userRepo: userRepo}
}

func (svc *UserSvc) CreateUser(name, email, password string) (*mdl.User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &mdl.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	if err := svc.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserSvc) GetUserByID(id uint) (*mdl.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserSvc) GetUserByEmail(email string) (*mdl.User, error) {
    return s.userRepo.FindByEmail(email)
}