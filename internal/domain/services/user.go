// File: internal/domain/services/user.go
// Tạo tại: internal/domain/services/user.go
// Mục đích: Service xử lý logic nghiệp vụ cho User (CRUD, assign roles, validate)

package services

import (
	"errors"
	"math"

	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/dto/request"
	"github.com/godiidev/appsynex/internal/dto/response"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUsers(req request.UserFilterRequest) (*response.PaginatedResponse, error)
	GetUserByID(id uint) (*response.UserDetailResponse, error)
	CreateUser(req request.CreateUserRequest) (*response.UserDetailResponse, error)
	UpdateUser(id uint, req request.UpdateUserRequest) (*response.UserDetailResponse, error)
	DeleteUser(id uint) error
	AssignRoles(userID uint, roleIDs []uint) (*response.UserDetailResponse, error)
}

type userService struct {
	userRepo interfaces.UserRepository
	roleRepo interfaces.RoleRepository
}

func NewUserService(userRepo interfaces.UserRepository, roleRepo interfaces.RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *userService) GetUsers(req request.UserFilterRequest) (*response.PaginatedResponse, error) {
	// Set defaults for pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Get users from repository
	users, total, err := s.userRepo.FindAll(req.Page, req.Limit, req.Search)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	items := make([]interface{}, len(users))
	for i, user := range users {
		items[i] = convertUserToResponse(&user)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &response.PaginatedResponse{
		Items:      items,
		TotalItems: total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) GetUserByID(id uint) (*response.UserDetailResponse, error) {
	user, err := s.userRepo.FindByIDWithRoles(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return convertUserToDetailResponse(user), nil
}

func (s *userService) CreateUser(req request.CreateUserRequest) (*response.UserDetailResponse, error) {
	// Check if username already exists
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:      req.Username,
		PasswordHash:  string(passwordHash),
		Email:         req.Email,
		Phone:         req.Phone,
		AccountStatus: "active",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Assign roles if provided
	if len(req.RoleIDs) > 0 {
		if err := s.userRepo.AssignRoles(user.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	// Get complete user with roles
	createdUser, err := s.userRepo.FindByIDWithRoles(user.ID)
	if err != nil {
		return nil, err
	}

	return convertUserToDetailResponse(createdUser), nil
}

func (s *userService) UpdateUser(id uint, req request.UpdateUserRequest) (*response.UserDetailResponse, error) {
	// Get existing user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if username is being changed and already exists
	if req.Username != "" && req.Username != user.Username {
		existingUser, _ := s.userRepo.FindByUsername(req.Username)
		if existingUser != nil {
			return nil, errors.New("username already exists")
		}
	}

	// Update fields if provided
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.AccountStatus != "" {
		user.AccountStatus = req.AccountStatus
	}
	if req.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(passwordHash)
	}

	// Update user
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Update roles if provided
	if len(req.RoleIDs) > 0 {
		if err := s.userRepo.AssignRoles(user.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	// Get updated user with roles
	updatedUser, err := s.userRepo.FindByIDWithRoles(id)
	if err != nil {
		return nil, err
	}

	return convertUserToDetailResponse(updatedUser), nil
}

func (s *userService) DeleteUser(id uint) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

func (s *userService) AssignRoles(userID uint, roleIDs []uint) (*response.UserDetailResponse, error) {
	// Check if user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate role IDs
	for _, roleID := range roleIDs {
		_, err := s.roleRepo.FindByID(roleID)
		if err != nil {
			return nil, errors.New("one or more roles not found")
		}
	}

	// Assign roles
	if err := s.userRepo.AssignRoles(userID, roleIDs); err != nil {
		return nil, err
	}

	// Get updated user with roles
	user, err := s.userRepo.FindByIDWithRoles(userID)
	if err != nil {
		return nil, err
	}

	return convertUserToDetailResponse(user), nil
}

// Helper functions to convert model to response DTO
func convertUserToResponse(user *models.User) *response.UserResponse {
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.RoleName
	}

	return &response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
	}
}

func convertUserToDetailResponse(user *models.User) *response.UserDetailResponse {
	roles := make([]response.RoleResponse, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = response.RoleResponse{
			ID:          role.ID,
			RoleName:    role.RoleName,
			Description: role.Description,
		}
	}

	return &response.UserDetailResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Phone:         user.Phone,
		LastLogin:     user.LastLogin,
		AccountStatus: user.AccountStatus,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Roles:         roles,
	}
}