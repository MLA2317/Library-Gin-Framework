package repository

import (
    "database/sql"
    "library-project/internal/models"
    "github.com/google/uuid"
)

type CategoryRepository struct {
    db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
    return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
    category.ID = uuid.New().String()
    
    query := `
        INSERT INTO categories (id, name, description)
        VALUES ($1, $2, $3)
        RETURNING created_at, updated_at
    `
    
    return r.db.QueryRow(query, category.ID, category.Name, category.Description).
        Scan(&category.CreatedAt, &category.UpdatedAt)
}

func (r *CategoryRepository) FindAll() ([]*models.Category, error) {
    query := `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var categories []*models.Category
    for rows.Next() {
        cat := &models.Category{}
        err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)
        if err != nil {
            return nil, err
        }
        categories = append(categories, cat)
    }
    
    return categories, nil
}

func (r *CategoryRepository) FindByID(id string) (*models.Category, error) {
    category := &models.Category{}
    
    query := `SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`
    
    err := r.db.QueryRow(query, id).Scan(
        &category.ID, &category.Name, &category.Description,
        &category.CreatedAt, &category.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    
    return category, err
}