package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (phone, name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Phone, user.Name).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, phone, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Phone, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByPhone retrieves a user by phone number
func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `
		SELECT id, phone, name, created_at, updated_at
		FROM users
		WHERE phone = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, phone).
		Scan(&user.ID, &user.Phone, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET name = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query, user.ID, user.Name).Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrNotFound
		}
		return err
	}

	return nil
}

// GetOrCreate gets a user by phone or creates a new one
func (r *UserRepository) GetOrCreate(ctx context.Context, phone string) (*models.User, bool, error) {
	// Try to get existing user
	user, err := r.GetByPhone(ctx, phone)
	if err == nil {
		return user, false, nil // User exists
	}
	if !errors.Is(err, models.ErrNotFound) {
		return nil, false, err // Other error
	}

	// Create new user
	user = &models.User{Phone: phone}
	if err := r.Create(ctx, user); err != nil {
		return nil, false, err
	}

	return user, true, nil // New user created
}
