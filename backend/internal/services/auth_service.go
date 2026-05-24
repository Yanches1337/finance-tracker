package services

import (
	"backend/internal/adapters/interfaces"
	"backend/internal/domain"
	"backend/internal/utils"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     interfaces.UserRepository
	tokenService *TokenService
	Cfg          *utils.Config
}

func NewAuthService(userRepo interfaces.UserRepository, tokenService *TokenService, cfg *utils.Config) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
		Cfg:          cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest, ip, userAgent string) (string, string, *domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", "", nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", "", nil, errors.New("invalid email or password")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.Cfg.JWT.SecretKey, s.Cfg.JWT.AccessTokenTTL)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.Cfg.JWT.SecretKey, s.Cfg.JWT.RefreshTokenTTL)
	if err != nil {
		return "", "", nil, err
	}

	if err := s.tokenService.SaveRefreshToken(ctx, user.ID, refreshToken, ip, userAgent); err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenService.RevokeRefreshToken(ctx, refreshToken)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	userID, err := s.tokenService.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	// Новый Access Token
	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.Cfg.JWT.SecretKey, s.Cfg.JWT.AccessTokenTTL)
	if err != nil {
		return "", "", err
	}

	// Новый Refresh Token (ротация)
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, s.Cfg.JWT.SecretKey, s.Cfg.JWT.RefreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	// Удаляем старый refresh token
	s.tokenService.RevokeRefreshToken(ctx, refreshToken)

	// Сохраняем новый
	s.tokenService.SaveRefreshToken(ctx, user.ID, newRefreshToken, "", "")

	return newAccessToken, newRefreshToken, nil
}
