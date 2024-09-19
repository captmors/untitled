package repo

import (
	"untitled/internal/users/mdl"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *mdl.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) FindByID(id uint) (*mdl.User, error) {
	var user mdl.User
	err := r.db.First(&user, id).Error
	return &user, err
}
