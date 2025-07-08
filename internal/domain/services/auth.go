package services

import (
	"errors"
	"log"
	"time"

	"github.com/godiidev/appsynex/internal/dto/request"
	"github.com/godiidev/appsynex/internal/dto/response"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
	"github.com/godiidev/appsynex/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req request.LoginRequest) (*response.LoginResponse, error)
}

type authService struct {
	userRepo       interfaces.UserRepository
	roleRepo       interfaces.RoleRepository
	permissionRepo interfaces.PermissionRepository
	jwtService     auth.JWTService
}

func NewAuthService(userRepo interfaces.UserRepository, roleRepo interfaces.RoleRepository, permissionRepo interfaces.PermissionRepository, jwtService auth.JWTService) AuthService {
	return &authService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		jwtService:     jwtService,
	}
}

func (s *authService) Login(req request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Get user with roles
	userWithRoles, err := s.userRepo.FindByIDWithRoles(user.ID)
	if err != nil {
		return nil, err
	}

	// Create roles list
	roles := make([]string, 0)
	var permissions []auth.Permission

	for _, role := range userWithRoles.Roles {
		roles = append(roles, role.RoleName)
	}

	// Get user's effective permissions if permission repo is available
	if s.permissionRepo != nil {
		effectivePerms, err := s.permissionRepo.GetUserEffectivePermissions(user.ID)
		if err == nil {
			for _, p := range effectivePerms {
				permissions = append(permissions, auth.Permission{
					Name:   p.PermissionName,
					Module: p.Module,
				})
			}
		} else {
			log.Printf("Warning: Could not get effective permissions for user %d: %v", user.ID, err)
		}
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user.ID, user.Username, roles, permissions)
	if err != nil {
		return nil, err
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := s.userRepo.Update(user); err != nil {
		// Just log this error, don't fail the login
		log.Printf("Failed to update last login: %v", err)
	}

	return &response.LoginResponse{
		Token: token,
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Roles:    roles,
		},
	}, nil
}
