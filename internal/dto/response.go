package dto

import (
	"library-project/internal/models"
	"time"
)

// Generic Response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Auth Responses
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string            `json:"id"`
	Email     string            `json:"email"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Role      models.UserRole   `json:"role"`
	CreatedAt time.Time         `json:"created_at"`
}

// Book Responses
type BookResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	PDFFile      string    `json:"pdf_file"`
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	OwnerID      string    `json:"owner_id"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
	SaveCount    int       `json:"save_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BookListResponse struct {
	Books      []BookResponse     `json:"books"`
	Pagination PaginationResponse `json:"pagination"`
}

type BookDetailResponse struct {
	BookResponse
	Owner    UserResponse      `json:"owner"`
}

// Category Responses
type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BookCount   int       `json:"book_count,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

// Comment Responses
type CommentResponse struct {
	ID            string    `json:"id"`
	BookID        string    `json:"book_id"`
	UserID        string    `json:"user_id"`
	UserFirstName string    `json:"user_first_name"`
	UserLastName  string    `json:"user_last_name"`
	Content       string    `json:"content"`
	ParentID      *string   `json:"parent_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
	Total    int               `json:"total"`
}

// Like Responses
type LikeResponse struct {
	BookID       string `json:"book_id"`
	UserLiked    bool   `json:"user_liked"`
	UserDisliked bool   `json:"user_disliked"`
	LikeCount    int    `json:"like_count"`
	DislikeCount int    `json:"dislike_count"`
}

// Saved Book Responses
type SavedBookResponse struct {
	ID        string       `json:"id"`
	Book      BookResponse `json:"book"`
	SavedAt   time.Time    `json:"saved_at"`
}

type SavedBookListResponse struct {
	Books      []SavedBookResponse `json:"books"`
	Total      int                 `json:"total"`
	Pagination PaginationResponse  `json:"pagination"`
}

// Pagination Response
type PaginationResponse struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// Dashboard Responses
type DashboardResponse struct {
	TotalBooks      int            `json:"total_books"`
	TotalCategories int            `json:"total_categories"`
	TotalSaves      int            `json:"total_saves,omitempty"`
	TotalComments   int            `json:"total_comments,omitempty"`
	PopularBooks    []BookResponse `json:"popular_books"`
	RecentBooks     []BookResponse `json:"recent_books"`
}

// Statistics Response
type BookStatisticsResponse struct {
	BookID       string `json:"book_id"`
	Views        int    `json:"views"`
	Saves        int    `json:"saves"`
	Likes        int    `json:"likes"`
	Dislikes     int    `json:"dislikes"`
	Comments     int    `json:"comments"`
}