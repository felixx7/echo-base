package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"echo-base/domain/entity"
	"echo-base/domain/usecase"
	"echo-base/utils"
)

// UserHandler handles user HTTP requests
type UserHandler struct {
	userUsecase usecase.UserUsecase
	validator   *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		validator:   validator.New(),
	}
}

// Register handles user registration
// POST /api/auth/register
func (h *UserHandler) Register(c echo.Context) error {
	payload := new(entity.UserCreatePayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request body"))
	}

	// Validate payload
	if err := h.validator.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
	}

	result, err := h.userUsecase.Register(payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse("user registered successfully", result))
}

// Login handles user login
// POST /api/auth/login
func (h *UserHandler) Login(c echo.Context) error {
	payload := new(entity.UserLoginPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request body"))
	}

	// Validate payload
	if err := h.validator.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
	}

	result, err := h.userUsecase.Login(payload)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("login successful", result))
}

// GetByID gets user by ID
// GET /api/users/:id
func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid user ID"))
	}

	result, err := h.userUsecase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("user retrieved successfully", result))
}

// GetAll gets all users
// GET /api/users
func (h *UserHandler) GetAll(c echo.Context) error {
	result, err := h.userUsecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("users retrieved successfully", result))
}

// GetAllPagination gets all users with pagination and optional search
// GET /api/users/pagination?page=1&limit=10&search=john
func (h *UserHandler) GetAllPagination(c echo.Context) error {
	// Parse pagination parameters
	page := int64(1)
	limit := int64(10)
	search := ""

	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.ParseInt(p, 10, 64); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 64); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if s := c.QueryParam("search"); s != "" {
		search = s
	}

	result, err := h.userUsecase.GetAllPagination(page, limit, search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("users retrieved successfully", result))
}

// Update updates user profile
// PUT /api/users/:id
func (h *UserHandler) Update(c echo.Context) error {
	// Check authorization
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("unauthorized"))
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid user ID"))
	}

	// Check if user is updating their own profile
	if userID.(int64) != id {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("you can only update your own profile"))
	}

	payload := new(struct {
		Name string `json:"name" validate:"required,min=3"`
	})
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request body"))
	}

	if err := h.validator.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
	}

	result, err := h.userUsecase.Update(id, payload.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("user updated successfully", result))
}

// Delete deletes a user
// DELETE /api/users/:id
func (h *UserHandler) Delete(c echo.Context) error {
	// Check authorization
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("unauthorized"))
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid user ID"))
	}

	// Check if user is deleting their own account
	if userID.(int64) != id {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("you can only delete your own account"))
	}

	err = h.userUsecase.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("user deleted successfully", nil))
}

// GetProfile gets current user profile
// GET /api/profile
func (h *UserHandler) GetProfile(c echo.Context) error {
	// Check authorization
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("unauthorized"))
	}

	result, err := h.userUsecase.GetByID(userID.(int64))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("profile retrieved successfully", result))
}
