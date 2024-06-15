package users

import (
	"context"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/users"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	log     *zap.Logger
	storage users.Storage
}

func New(log *zap.Logger, storage users.Storage) Service {
	return &service{log: log, storage: storage}
}

func (s *service) Register(ctx context.Context, userIn models.User) (*models.User, error) {
	hashPassword, err := s.getHashPassword(userIn.Password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Login:    userIn.Login,
		Password: hashPassword,
	}
	registeredUser, err := s.storage.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	return registeredUser, nil
}

func (s *service) Login(ctx context.Context, userIn models.User) (*models.User, error) {
	user, err := s.storage.Login(ctx, userIn.Login)
	if err != nil {
		return nil, err
	}

	if !s.checkPassword(user.Password, userIn.Password) {
		return nil, ErrIncorrectData
	}
	return user, nil
}

func (s *service) getHashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *service) checkPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
