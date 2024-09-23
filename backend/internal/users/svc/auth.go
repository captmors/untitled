package svc

import (
	"errors"
	"strconv"
	"time"
	"untitled/internal/users/mdl"
	"untitled/internal/users/repo"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthSvc struct {
	userRepo *repo.UserRepo
	jwtKey   []byte
}

func NewAuthSvc(userRepo *repo.UserRepo, jwtKey []byte) *AuthSvc {
	return &AuthSvc{
		userRepo: userRepo,
		jwtKey:   jwtKey,
	}
}

func (a *AuthSvc) GenerateToken(user *mdl.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Subject:   strconv.Itoa(int(user.ID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtKey)
}

func (a *AuthSvc) AuthenticateByName(name, password string) (*mdl.User, error) {
	user, err := a.userRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	if err := checkPasswordHash(password, user.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (a *AuthSvc) ValidateToken(tokenStr string) (*mdl.User, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, err
	}

	return a.userRepo.FindByID(uint(userID))
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (a *AuthSvc) Register(name, email, password string) (*mdl.User, error) {
	if _, err := a.userRepo.FindByEmail(email); err == nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &mdl.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	if err := a.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
