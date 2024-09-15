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

func (s *UserSvc) CreateUser(name, email, password string) (*mdl.User, error) {
	user := &mdl.User{Name: name, Email: email, Password: password}
	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserSvc) GetUserByID(id uint) (*mdl.User, error) {
	return s.userRepo.FindByID(id)
}

