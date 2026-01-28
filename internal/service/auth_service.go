package service

import (
    "errors"
    "library-project/config"
    "library-project/internal/dto"
    "library-project/internal/models"
    "library-project/internal/repository"
    "library-project/internal/utils"
)

type AuthService struct {
    userRepo *repository.UserRepository
    cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        cfg:      cfg,
    }
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
    existing, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, errors.New("email already exists")
    }

    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    userRole := req.Role
    if userRole == "" {
        userRole = models.RoleMember
    }

    user := &models.User{
        Email:     req.Email,
        Password:  hashedPassword,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Role:      userRole,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    token, err := utils.GenerateToken(
        user.ID,
        user.Email,
        string(user.Role),
        s.cfg.JWT.Secret,
        s.cfg.JWT.ExpirationHours,
    )
    if err != nil {
        return nil, err
    }

    return &dto.AuthResponse{
        Token: token,
        User:  utils.MapUserToResponse(user),  // ← Changed
    }, nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid email or password")
    }

    if !utils.CheckPassword(req.Password, user.Password) {
        return nil, errors.New("invalid email or password")
    }

    token, err := utils.GenerateToken(
        user.ID,
        user.Email,
        string(user.Role),
        s.cfg.JWT.Secret,
        s.cfg.JWT.ExpirationHours,
    )
    if err != nil {
        return nil, err
    }

    return &dto.AuthResponse{
        Token: token,
        User:  utils.MapUserToResponse(user),  // ← Changed
    }, nil
}