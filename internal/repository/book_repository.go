package repository

import (
    "database/sql"
    "library-project/internal/models"
    "github.com/google/uuid"
)

type BookRepository struct {
    db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
    return &BookRepository{db: db}
}

func (r *BookRepository) Create(book *models.Book) error {
    book.ID = uuid.New().String()
    
    query := `
        INSERT INTO books (id, title, description, pdf_file, category_id, owner_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING created_at, updated_at
    `
    
    return r.db.QueryRow(query, book.ID, book.Title, book.Description,
        book.PDFFile, book.CategoryID, book.OwnerID).
        Scan(&book.CreatedAt, &book.UpdatedAt)
}

func (r *BookRepository) Update(book *models.Book) error {
    query := `
        UPDATE books 
        SET title = $1, description = $2, category_id = $3
        WHERE id = $4
        RETURNING updated_at
    `
    
    return r.db.QueryRow(query, book.Title, book.Description,
        book.CategoryID, book.ID).Scan(&book.UpdatedAt)
}

func (r *BookRepository) Delete(id string) error {
    query := `DELETE FROM books WHERE id = $1`
    _, err := r.db.Exec(query, id)
    return err
}

func (r *BookRepository) FindByID(id string) (*models.BookWithCategory, error) {
    book := &models.BookWithCategory{}
    
    query := `
        SELECT b.id, b.title, b.description, b.pdf_file, b.category_id, b.owner_id,
               b.like_count, b.dislike_count, b.save_count, b.created_at, b.updated_at,
               c.name as category_name
        FROM books b
        JOIN categories c ON b.category_id = c.id
        WHERE b.id = $1
    `
    
    err := r.db.QueryRow(query, id).Scan(
        &book.ID, &book.Title, &book.Description, &book.PDFFile, &book.CategoryID,
        &book.OwnerID, &book.LikeCount, &book.DislikeCount, &book.SaveCount,
        &book.CreatedAt, &book.UpdatedAt, &book.CategoryName,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    
    return book, err
}

func (r *BookRepository) FindAll() ([]*models.BookWithCategory, error) {
    return r.FindAllPaginated(0, 0)
}

func (r *BookRepository) FindAllPaginated(limit, offset int) ([]*models.BookWithCategory, error) {
    query := `
        SELECT b.id, b.title, b.description, b.pdf_file, b.category_id, b.owner_id,
               b.like_count, b.dislike_count, b.save_count, b.created_at, b.updated_at,
               c.name as category_name
        FROM books b
        JOIN categories c ON b.category_id = c.id
        ORDER BY b.save_count DESC
    `

    // Add pagination if limit > 0
    if limit > 0 {
        query += ` LIMIT $1 OFFSET $2`
    }

    var rows *sql.Rows
    var err error

    if limit > 0 {
        rows, err = r.db.Query(query, limit, offset)
    } else {
        rows, err = r.db.Query(query)
    }
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var books []*models.BookWithCategory
    for rows.Next() {
        book := &models.BookWithCategory{}
        err := rows.Scan(
            &book.ID, &book.Title, &book.Description, &book.PDFFile, &book.CategoryID,
            &book.OwnerID, &book.LikeCount, &book.DislikeCount, &book.SaveCount,
            &book.CreatedAt, &book.UpdatedAt, &book.CategoryName,
        )
        if err != nil {
            return nil, err
        }
        books = append(books, book)
    }
    
    return books, nil
}

func (r *BookRepository) FindByCategory(categoryID string) ([]*models.BookWithCategory, error) {
    return r.FindByCategoryPaginated(categoryID, 0, 0)
}

func (r *BookRepository) FindByCategoryPaginated(categoryID string, limit, offset int) ([]*models.BookWithCategory, error) {
    query := `
        SELECT b.id, b.title, b.description, b.pdf_file, b.category_id, b.owner_id,
               b.like_count, b.dislike_count, b.save_count, b.created_at, b.updated_at,
               c.name as category_name
        FROM books b
        JOIN categories c ON b.category_id = c.id
        WHERE b.category_id = $1
        ORDER BY b.save_count DESC
    `

    // Add pagination if limit > 0
    if limit > 0 {
        query += ` LIMIT $2 OFFSET $3`
    }

    var rows *sql.Rows
    var err error

    if limit > 0 {
        rows, err = r.db.Query(query, categoryID, limit, offset)
    } else {
        rows, err = r.db.Query(query, categoryID)
    }
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var books []*models.BookWithCategory
    for rows.Next() {
        book := &models.BookWithCategory{}
        err := rows.Scan(
            &book.ID, &book.Title, &book.Description, &book.PDFFile, &book.CategoryID,
            &book.OwnerID, &book.LikeCount, &book.DislikeCount, &book.SaveCount,
            &book.CreatedAt, &book.UpdatedAt, &book.CategoryName,
        )
        if err != nil {
            return nil, err
        }
        books = append(books, book)
    }
    
    return books, nil
}

func (r *BookRepository) UpdateLikeCount(bookID string) error {
    query := `
        UPDATE books 
        SET like_count = (SELECT COUNT(*) FROM likes WHERE book_id = $1 AND is_like = true),
            dislike_count = (SELECT COUNT(*) FROM likes WHERE book_id = $1 AND is_like = false)
        WHERE id = $1
    `
    _, err := r.db.Exec(query, bookID)
    return err
}

func (r *BookRepository) UpdateSaveCount(bookID string) error {
    query := `
        UPDATE books
        SET save_count = (SELECT COUNT(*) FROM saved_books WHERE book_id = $1)
        WHERE id = $1
    `
    _, err := r.db.Exec(query, bookID)
    return err
}

func (r *BookRepository) CountAll() (int, error) {
    var count int
    query := `SELECT COUNT(*) FROM books`
    err := r.db.QueryRow(query).Scan(&count)
    return count, err
}

func (r *BookRepository) CountByCategory(categoryID string) (int, error) {
    var count int
    query := `SELECT COUNT(*) FROM books WHERE category_id = $1`
    err := r.db.QueryRow(query, categoryID).Scan(&count)
    return count, err
}