package repository

import (
    "database/sql"
    "library-project/internal/models"
    "github.com/google/uuid"
)

type CommentRepository struct {
    db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
    return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
    comment.ID = uuid.New().String()
    
    query := `
        INSERT INTO comments (id, book_id, user_id, content, parent_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING created_at, updated_at
    `
    
    return r.db.QueryRow(query, comment.ID, comment.BookID, comment.UserID,
        comment.Content, comment.ParentID).Scan(&comment.CreatedAt, &comment.UpdatedAt)
}

func (r *CommentRepository) FindByBookID(bookID string) ([]*models.CommentWithUser, error) {
    query := `
        SELECT c.id, c.book_id, c.user_id, c.content, c.parent_id, c.created_at, c.updated_at,
               u.first_name, u.last_name
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.book_id = $1
        ORDER BY c.created_at DESC
    `
    
    rows, err := r.db.Query(query, bookID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var comments []*models.CommentWithUser
    for rows.Next() {
        comment := &models.CommentWithUser{}
        err := rows.Scan(
            &comment.ID, &comment.BookID, &comment.UserID, &comment.Content,
            &comment.ParentID, &comment.CreatedAt, &comment.UpdatedAt,
            &comment.UserFirstName, &comment.UserLastName,
        )
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }
    
    return comments, nil
}

type LikeRepository struct {
    db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
    return &LikeRepository{db: db}
}

func (r *LikeRepository) Upsert(like *models.Like) error {
    like.ID = uuid.New().String()
    
    query := `
        INSERT INTO likes (id, book_id, user_id, is_like)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, book_id)
        DO UPDATE SET is_like = $4
        RETURNING created_at
    `
    
    return r.db.QueryRow(query, like.ID, like.BookID, like.UserID, like.IsLike).
        Scan(&like.CreatedAt)
}

func (r *LikeRepository) Delete(userID, bookID string) error {
    query := `DELETE FROM likes WHERE user_id = $1 AND book_id = $2`
    _, err := r.db.Exec(query, userID, bookID)
    return err
}

type SavedBookRepository struct {
    db *sql.DB
}

func NewSavedBookRepository(db *sql.DB) *SavedBookRepository {
    return &SavedBookRepository{db: db}
}

func (r *SavedBookRepository) Create(saved *models.SavedBook) error {
    saved.ID = uuid.New().String()
    
    query := `
        INSERT INTO saved_books (id, user_id, book_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, book_id) DO NOTHING
        RETURNING created_at
    `
    
    return r.db.QueryRow(query, saved.ID, saved.UserID, saved.BookID).
        Scan(&saved.CreatedAt)
}

func (r *SavedBookRepository) Delete(userID, bookID string) error {
    query := `DELETE FROM saved_books WHERE user_id = $1 AND book_id = $2`
    _, err := r.db.Exec(query, userID, bookID)
    return err
}

func (r *SavedBookRepository) FindByUserID(userID string) ([]*models.BookWithCategory, error) {
    query := `
        SELECT b.id, b.title, b.description, b.pdf_file, b.category_id, b.owner_id,
               b.like_count, b.dislike_count, b.save_count, b.created_at, b.updated_at,
               c.name as category_name
        FROM saved_books sb
        JOIN books b ON sb.book_id = b.id
        JOIN categories c ON b.category_id = c.id
        WHERE sb.user_id = $1
        ORDER BY sb.created_at DESC
    `
    
    rows, err := r.db.Query(query, userID)
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