package services

import (
	"errors"
	"log"
	"time"

	"github.com/godiidev/appsynex/internal/domain/models"
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
	userRepo   interfaces.UserRepository
	roleRepo   interfaces.RoleRepository
	jwtService auth.JWTService
}

func NewAuthService(userRepo interfaces.UserRepository, roleRepo interfaces.RoleRepository, jwtService auth.JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		jwtService: jwtService,
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

	// Create roles and permissions lists
	roles := make([]string, 0)
	var permissions []auth.Permission

	for _, role := range userWithRoles.Roles {
		roles = append(roles, role.RoleName)

		// Get role with permissions
		roleWithPermissions, err := s.roleRepo.FindByIDWithPermissions(role.ID)
		if err != nil {
			continue
		}

		for _, p := range roleWithPermissions.Permissions {
			// Find module from role_permission
			var rp models.RolePermission
			if err := s.roleRepo.GetDB().Where("role_id = ? AND permission_id = ?", role.ID, p.ID).First(&rp).Error; err != nil {
				continue
			}

			permissions = append(permissions, auth.Permission{
				Name:   p.PermissionName,
				Module: rp.Module,
			})
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
