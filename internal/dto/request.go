package dto

// Auth Requests
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Book Requests
type CreateBookRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description" binding:"required"`
	CategoryID  string `form:"category_id" binding:"required"`
}

type UpdateBookRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	CategoryID  string `json:"category_id" binding:"required"`
}

// Category Requests
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// Comment Requests
type CreateCommentRequest struct {
	Content  string  `json:"content" binding:"required,min=1,max=1000"`
	ParentID *string `json:"parent_id,omitempty"`
}

// Like Requests
type LikeRequest struct {
	IsLike bool `json:"is_like"`
}

type RemoveLikeRequest struct {
	BookID string `json:"book_id" binding:"required"`
}

// Pagination Request
type PaginationRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy   string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at save_count title"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// Filter Request
type BookFilterRequest struct {
	CategoryID string `form:"category_id"`
	Search     string `form:"search"`
	PaginationRequest
}