package svc

import (
	"errors"
	"strconv"
	"time"
	"untitled/internal/users/mdl"
	"untitled/internal/users/repo"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo *repo.UserRepo
    jwtKey   []byte
}

func NewAuthService(userRepo *repo.UserRepo, jwtKey []byte) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        jwtKey:   jwtKey,
    }
}

func (a *AuthService) GenerateToken(user *mdl.User) (string, error) {
    claims := &jwt.StandardClaims{
        ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        Subject:   strconv.Itoa(int(user.ID)),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(a.jwtKey)
}


func (a *AuthService) Authenticate(email, password string) (*mdl.User, error) {
    user, err := a.userRepo.FindByEmail(email)
    if err != nil {
        return nil, err
    }

    if err := checkPasswordHash(password, user.Password); err != nil {
        return nil, errors.New("invalid credentials")
    }

    return user, nil
}


func (a *AuthService) ValidateToken(tokenStr string) (*mdl.User, error) {
    claims := &jwt.StandardClaims{}
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