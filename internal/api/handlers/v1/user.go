// File: internal/api/handlers/v1/user.go
// Tạo tại: internal/api/handlers/v1/user.go
// Mục đích: Handler xử lý các API liên quan đến quản lý người dùng (CRUD users, assign roles)

package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/godiidev/appsynex/internal/domain/services"
	"github.com/godiidev/appsynex/internal/dto/request"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAll godoc
// @Summary     Get all users
// @Description Get a list of users with pagination and filtering
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       page query int false "Page number"
// @Param       limit query int false "Items per page"
// @Param       search query string false "Search term for username or email"
// @Security    BearerAuth
// @Success     200 {object} response.PaginatedResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	var req request.UserFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.userService.GetUsers(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetByID godoc
// @Summary     Get user by ID
// @Description Get a user by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Security    BearerAuth
// @Success     200 {object} response.UserDetailResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Create godoc
// @Summary     Create a new user
// @Description Create a new user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       user body request.CreateUserRequest true "User to create"
// @Security    BearerAuth
// @Success     201 {object} response.UserDetailResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Update godoc
// @Summary     Update a user
// @Description Update a user by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Param       user body request.UpdateUserRequest true "User data to update"
// @Security    BearerAuth
// @Success     200 {object} response.UserDetailResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete godoc
// @Summary     Delete a user
// @Description Delete a user by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Security    BearerAuth
// @Success     204 {object} nil
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignRoles godoc
// @Summary     Assign roles to user
// @Description Assign roles to a user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Param       roles body request.AssignRolesRequest true "Roles to assign"
// @Security    BearerAuth
// @Success     200 {object} response.UserDetailResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /users/{id}/roles [post]
func (h *UserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req request.AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.AssignRoles(uint(id), req.RoleIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}