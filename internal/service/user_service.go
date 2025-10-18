package service

import (
	"errors"
	"httpServerTest/internal/database"
	"httpServerTest/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(username, password string) error {
	var count int64
	database.DB.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Username: username,
		Password: string(hash),
	}
	return database.DB.Create(&user).Error
}

func (s *UserService) Login(username, password string) (model.User, error) {
	var user model.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if username == "hwt" {
			var userCount int64
			database.DB.Model(&model.User{}).Count(&userCount)
			if userCount == 0 {
				err := s.Register(username, password)
				if err != nil {
					return model.User{}, err
				}
				return model.User{}, errors.New("初始用户注册成功，请重新登录！")
			}
		} else {
			return model.User{}, errors.New("用户不存在")
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.User{}, errors.New("密码错误")
	}
	return user, nil
}
