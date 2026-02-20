package usecase

import (
	"errors"
	"fmt"

	"echo-base/domain/entity"
	"echo-base/domain/repository"
	"echo-base/utils"
)

// UserUsecase defines the interface for user usecase
type UserUsecase interface {
	// Register registers a new user
	Register(payload *entity.UserCreatePayload) (*entity.UserResponse, error)

	// Login logs in a user and returns a token
	Login(payload *entity.UserLoginPayload) (*entity.LoginResponse, error)

	// GetByID gets a user by ID
	GetByID(id int64) (*entity.UserResponse, error)

	// GetAll gets all users
	GetAll() ([]*entity.UserResponse, error)

	// GetAllPagination gets all users with pagination and optional search
	GetAllPagination(page int64, limit int64, search string) (*entity.PaginatedUserResponse, error)

	// Update updates a user
	Update(id int64, name string) (*entity.UserResponse, error)

	// Delete deletes a user
	Delete(id int64) error
}

// UserUsecaseImpl implements UserUsecase
type UserUsecaseImpl struct {
	userRepo repository.UserRepository
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
	}
}

// Register registers a new user
func (u *UserUsecaseImpl) Register(payload *entity.UserCreatePayload) (*entity.UserResponse, error) {
	// Check if email is already registered
	existingUser, err := u.userRepo.GetByEmail(payload.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email is already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create user
	user := &entity.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	createdUser, err := u.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &entity.UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}

// Login logs in a user and returns a token
func (u *UserUsecaseImpl) Login(payload *entity.UserLoginPayload) (*entity.LoginResponse, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(payload.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !utils.CheckPassword(user.Password, payload.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token with role
	token, err := utils.GenerateToken(user.ID, user.Email, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &entity.LoginResponse{
		Token: token,
		User: entity.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			RoleID:    user.RoleID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

// GetByID gets a user by ID
func (u *UserUsecaseImpl) GetByID(id int64) (*entity.UserResponse, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &entity.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetAll gets all users
func (u *UserUsecaseImpl) GetAll() ([]*entity.UserResponse, error) {
	users, err := u.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}

	responses := make([]*entity.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, &entity.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			RoleID:    user.RoleID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return responses, nil
}

// Update updates a user
func (u *UserUsecaseImpl) Update(id int64, name string) (*entity.UserResponse, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Name = name

	updatedUser, err := u.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &entity.UserResponse{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		RoleID:    updatedUser.RoleID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

// Delete deletes a user
func (u *UserUsecaseImpl) Delete(id int64) error {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	return u.userRepo.Delete(id)
}

// GetAllPagination gets all users with pagination and optional search
func (u *UserUsecaseImpl) GetAllPagination(page int64, limit int64, search string) (*entity.PaginatedUserResponse, error) {
	users, total, err := u.userRepo.GetAllPagination(page, limit, search)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}

	responses := make([]*entity.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, &entity.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			RoleID:    user.RoleID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	return &entity.PaginatedUserResponse{
		Data: responses,
		Pagination: entity.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
