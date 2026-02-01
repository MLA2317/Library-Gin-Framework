package service

import (
    "library-project/config"
    "library-project/internal/dto"
    "library-project/internal/models"
    "library-project/internal/repository"
    "library-project/internal/utils"
    "time"
)

type AuthService struct {
    userRepo         *repository.UserRepository
    refreshTokenRepo *repository.RefreshTokenRepository
    cfg              *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, refreshTokenRepo *repository.RefreshTokenRepository, cfg *config.Config) *AuthService {
    return &AuthService{
        userRepo:         userRepo,
        refreshTokenRepo: refreshTokenRepo,
        cfg:              cfg,
    }
}

func (s *AuthService) generateTokenPair(user *models.User) (*dto.AuthResponse, error) {
    accessToken, err := utils.GenerateToken(
        user.ID,
        user.Email,
        string(user.Role),
        s.cfg.JWT.Secret,
        s.cfg.JWT.ExpirationHours,
    )
    if err != nil {
        return nil, utils.NewInternalServerError("failed to generate access token", err)
    }

    refreshTokenStr, err := utils.GenerateRefreshToken()
    if err != nil {
        return nil, utils.NewInternalServerError("failed to generate refresh token", err)
    }

    refreshToken := &models.RefreshToken{
        UserID:    user.ID,
        Token:     refreshTokenStr,
        ExpiresAt: time.Now().AddDate(0, 0, s.cfg.JWT.RefreshExpirationDays),
    }

    if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
        return nil, utils.NewInternalServerError("failed to save refresh token", err)
    }

    return &dto.AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshTokenStr,
        ExpiresIn:    s.cfg.JWT.ExpirationHours * 3600,
        User:         utils.MapUserToResponse(user),
    }, nil
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
    if err := utils.ValidateEmail(req.Email); err != nil {
        return nil, utils.NewValidationError(err.Error())
    }

    if err := utils.ValidatePassword(req.Password); err != nil {
        return nil, utils.NewValidationError(err.Error())
    }

    existing, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, utils.NewInternalServerError("failed to check existing user", err)
    }
    if existing != nil {
        return nil, utils.NewAlreadyExistsError("user with this email")
    }

    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, utils.NewInternalServerError("failed to hash password", err)
    }

    user := &models.User{
        Email:     req.Email,
        Password:  hashedPassword,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Role:      models.RoleMember,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    return s.generateTokenPair(user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, utils.NewInternalServerError("failed to find user", err)
    }
    if user == nil {
        return nil, utils.NewUnauthorizedError("invalid email or password")
    }

    if !utils.CheckPassword(req.Password, user.Password) {
        return nil, utils.NewUnauthorizedError("invalid email or password")
    }

    return s.generateTokenPair(user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
    rt, err := s.refreshTokenRepo.FindByToken(req.RefreshToken)
    if err != nil {
        return nil, utils.NewInternalServerError("failed to find refresh token", err)
    }
    if rt == nil {
        return nil, utils.NewUnauthorizedError("invalid refresh token")
    }

    if time.Now().After(rt.ExpiresAt) {
        s.refreshTokenRepo.DeleteByToken(req.RefreshToken)
        return nil, utils.NewUnauthorizedError("refresh token expired")
    }

    // Delete old refresh token
    s.refreshTokenRepo.DeleteByToken(req.RefreshToken)

    // Invalidate all old access tokens
    if err := s.userRepo.InvalidateTokens(rt.UserID); err != nil {
        return nil, utils.NewInternalServerError("failed to invalidate tokens", err)
    }

    user, err := s.userRepo.FindByID(rt.UserID)
    if err != nil {
        return nil, utils.NewInternalServerError("failed to find user", err)
    }
    if user == nil {
        return nil, utils.NewUnauthorizedError("user not found")
    }

    return s.generateTokenPair(user)
}

func (s *AuthService) Logout(userID string) error {
    // Invalidate all access tokens
    if err := s.userRepo.InvalidateTokens(userID); err != nil {
        return err
    }
    // Delete all refresh tokens
    return s.refreshTokenRepo.DeleteByUserID(userID)
}
