package mdl

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int64  `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
