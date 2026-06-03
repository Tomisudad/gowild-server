package service

import (
	"errors"

	"github.com/gowild/server/internal/model"
	"github.com/gowild/server/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists      = errors.New("用户已存在")
	ErrUserNotFound    = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserService) Register(req model.RegisterRequest) (*model.User, error) {
	// 检查邮箱是否已注册
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return nil, ErrUserExists
	}

	// 检查用户名是否已存在
	if _, err := s.repo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已被占用")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req model.LoginRequest) (*model.User, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *UserService) GetProfile(userID int) (*model.User, error) {
	return s.repo.FindByID(userID)
}

func (s *UserService) UpdateProfile(userID int, req model.UpdateProfileRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	return s.repo.Update(user)
}

func (s *UserService) UpdateStats(userID int, dist float64) error {
	return s.repo.UpdateStats(userID, dist)
}

func (s *UserService) GetJWTSecret() string {
	return s.jwtSecret
}
