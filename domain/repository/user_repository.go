package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"echo-base/domain/entity"
)

// UserRepository defines the interface for user repository
type UserRepository interface {
	// GetByID gets a user by ID
	GetByID(id int64) (*entity.User, error)

	// GetByEmail gets a user by email
	GetByEmail(email string) (*entity.User, error)

	// Create creates a new user
	Create(user *entity.User) (*entity.User, error)

	// Update updates a user
	Update(user *entity.User) (*entity.User, error)

	// Delete deletes a user
	Delete(id int64) error

	// GetAll gets all users
	GetAll() ([]*entity.User, error)

	// GetAllPagination gets all users with pagination and optional search
	GetAllPagination(page int64, limit int64, search string) ([]*entity.User, int64, error)
}

// userRepository is a PostgreSQL implementation of UserRepository
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByID gets a user by ID from PostgreSQL
func (r *userRepository) GetByID(id int64) (*entity.User, error) {
	query := `
		SELECT id, name, email, password, role_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &entity.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}

	return user, nil
}

// GetByEmail gets a user by email from PostgreSQL
func (r *userRepository) GetByEmail(email string) (*entity.User, error) {
	query := `
		SELECT id, name, email, password, role_id, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &entity.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return user, nil
}

// Create creates a new user in PostgreSQL
func (r *userRepository) Create(user *entity.User) (*entity.User, error) {
	query := `
		INSERT INTO users (name, email, password, role_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Default to role_id 1 (user) if not set
	if user.RoleID == 0 {
		user.RoleID = 1
	}

	err := r.db.QueryRow(query,
		user.Name,
		user.Email,
		user.Password,
		user.RoleID,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

// Update updates a user in PostgreSQL
func (r *userRepository) Update(user *entity.User) (*entity.User, error) {
	query := `
		UPDATE users
		SET name = $1, role_id = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, name, email, password, role_id, created_at, updated_at
	`

	user.UpdatedAt = time.Now()

	err := r.db.QueryRow(query,
		user.Name,
		user.RoleID,
		user.UpdatedAt,
		user.ID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return user, nil
}

// Delete deletes a user from PostgreSQL
func (r *userRepository) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetAll gets all users from PostgreSQL
func (r *userRepository) GetAll() ([]*entity.User, error) {
	query := `
		SELECT id, name, email, password, role_id, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	users := make([]*entity.User, 0)
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.RoleID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return users, nil
}

// GetAllPagination gets all users with pagination and optional search
func (r *userRepository) GetAllPagination(page int64, limit int64, search string) ([]*entity.User, int64, error) {
	// Default pagination values
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Count total users
	countQuery := "SELECT COUNT(*) FROM users"
	var countArgs []interface{}

	if search != "" {
		countQuery += " WHERE name ILIKE $1 OR email ILIKE $1"
		countArgs = append(countArgs, "%"+search+"%")
	}

	var total int64
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting users: %w", err)
	}

	// Get data with pagination
	query := `
		SELECT id, name, email, password, role_id, created_at, updated_at
		FROM users
	`

	var args []interface{}
	argNum := 1

	if search != "" {
		query += fmt.Sprintf(" WHERE name ILIKE $%d OR email ILIKE $%d", argNum, argNum)
		args = append(args, "%"+search+"%")
		argNum++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	users := make([]*entity.User, 0, limit)
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.RoleID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error reading rows: %w", err)
	}

	return users, total, nil
}
