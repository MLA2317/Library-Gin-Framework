package repository

import (
    "database/sql"
    "library-project/internal/models"
    "github.com/google/uuid"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    user.ID = uuid.New().String()
    
    query := `
        INSERT INTO users (id, email, password, first_name, last_name, role)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING created_at, updated_at
    `
    
    return r.db.QueryRow(query, user.ID, user.Email, user.Password, 
        user.FirstName, user.LastName, user.Role).
        Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    user := &models.User{}
    
    query := `
        SELECT id, email, password, first_name, last_name, role, created_at, updated_at
        FROM users WHERE email = $1
    `
    
    err := r.db.QueryRow(query, email).Scan(
        &user.ID, &user.Email, &user.Password, &user.FirstName,
        &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    
    return user, err
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
    user := &models.User{}
    
    query := `
        SELECT id, email, password, first_name, last_name, role, created_at, updated_at
        FROM users WHERE id = $1
    `
    
    err := r.db.QueryRow(query, id).Scan(
        &user.ID, &user.Email, &user.Password, &user.FirstName,
        &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    
    return user, err
}