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

func (r *UserRepo) Create(user *mdl.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) FindByID(id uint) (*mdl.User, error) {
	var user mdl.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepo) FindByEmail(email string) (*mdl.User, error) {
	var user mdl.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepo) FindByName(name string) (*mdl.User, error) {
	var user mdl.User
	err := r.db.Where("name = ?", name).First(&user).Error
	return &user, err
}
