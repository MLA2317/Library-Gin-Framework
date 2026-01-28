package handler

import (
    "fmt"
    "library-project/config"
    "library-project/internal/dto"
    "library-project/internal/service"
    "net/http"
    "path/filepath"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type BookHandler struct {
    bookService *service.BookService
    cfg         *config.Config
}

func NewBookHandler(bookService *service.BookService, cfg *config.Config) *BookHandler {
    return &BookHandler{
        bookService: bookService,
        cfg:         cfg,
    }
}

// CreateBook godoc
// @Summary Create a new book (Owner only)
// @Description Upload a book with PDF file (Owner only)
// @Tags books
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Book Title"
// @Param description formData string true "Book Description"
// @Param category_id formData string true "Category ID"
// @Param pdf_file formData file true "PDF File"
// @Success 201 {object} dto.BookResponse 
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
    var req dto.CreateBookRequest
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    file, err := c.FormFile("pdf_file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "PDF file is required"})
        return
    }

    if file.Size > h.cfg.Upload.MaxFileSize {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds limit"})
        return
    }

    filename := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(file.Filename))
    filepath := filepath.Join(h.cfg.Upload.Path, filename)

    if err := c.SaveUploadedFile(file, filepath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    userID := c.GetString("user_id")
    book, err := h.bookService.CreateBook(&req, filename, userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, book)
}

// GetAllBooks godoc
// @Summary Get all books
// @Description Get all books ordered by save count
// @Tags books
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CategoryResponse
// @Failure 500 {object} map[string]string
// @Router /books [get]
func (h *BookHandler) GetAllBooks(c *gin.Context) {
    books, err := h.bookService.GetAllBooks()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, books)
}

// GetBook godoc
// @Summary Get a book by ID
// @Description Get detailed information about a book
// @Tags books
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {object} dto.BookResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
    id := c.Param("id")

    book, err := h.bookService.GetBook(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if book == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    c.JSON(http.StatusOK, book)
}

// UpdateBook godoc
// @Summary Update a book (Owner only)
// @Description Update book information (Owner only)
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Param request body dto.UpdateBookRequest true "Update Book Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
    id := c.Param("id")
    var req dto.UpdateBookRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.GetString("user_id")
    if err := h.bookService.UpdateBook(id, &req, userID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

// DeleteBook godoc
// @Summary Delete a book (Owner only)
// @Description Delete a book (Owner only)
// @Tags books
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
    id := c.Param("id")
    userID := c.GetString("user_id")

    if err := h.bookService.DeleteBook(id, userID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// GetBooksByCategory godoc
// @Summary Get books by category
// @Description Get all books in a specific category ordered by save count
// @Tags books
// @Produce json
// @Security BearerAuth
// @Param category_id query string true "Category ID"
// @Success 200 {array} dto.BookResponse
// @Failure 500 {object} map[string]string
// @Router /books/category [get]
func (h *BookHandler) GetBooksByCategory(c *gin.Context) {
    categoryID := c.Query("category_id")
    if categoryID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "category_id is required"})
        return
    }

    books, err := h.bookService.GetBooksByCategory(categoryID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, books)
}

// SaveBook godoc
// @Summary Save a book (Member only)
// @Description Save a book to user's collection (Member only)
// @Tags books
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /books/{id}/save [post]
func (h *BookHandler) SaveBook(c *gin.Context) {
    bookID := c.Param("id")
    userID := c.GetString("user_id")

    if err := h.bookService.SaveBook(userID, bookID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book saved successfully"})
}

// UnsaveBook godoc
// @Summary Unsave a book (Member only)
// @Description Remove a book from user's saved collection (Member only)
// @Tags books
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /books/{id}/unsave [delete]
func (h *BookHandler) UnsaveBook(c *gin.Context) {
    bookID := c.Param("id")
    userID := c.GetString("user_id")

    if err := h.bookService.UnsaveBook(userID, bookID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book unsaved successfully"})
}

// GetSavedBooks godoc
// @Summary Get saved books (Member only)
// @Description Get all books saved by the user (Member only)
// @Tags books
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.SavedBookResponse
// @Failure 500 {object} map[string]string
// @Router /books/saved [get]
func (h *BookHandler) GetSavedBooks(c *gin.Context) {
    userID := c.GetString("user_id")

    books, err := h.bookService.GetSavedBooks(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, books)
}

// LikeBook godoc
// @Summary Like or dislike a book
// @Description Add like or dislike to a book
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Param request body dto.LikeRequest true "Like Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /books/{id}/like [post]
func (h *BookHandler) LikeBook(c *gin.Context) {
    bookID := c.Param("id")
    userID := c.GetString("user_id")

    var req dto.LikeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.bookService.LikeBook(userID, bookID, req.IsLike); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Like updated successfully"})
}

// AddComment godoc
// @Summary Add a comment
// @Description Add a comment to a book (Member can comment, Owner can reply)
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Param request body dto.CreateCommentRequest true "Comment Request"
// @Success 201 {object} dto.CommentResponse
// @Failure 400 {object} map[string]string
// @Router /books/{id}/comments [post]
func (h *BookHandler) AddComment(c *gin.Context) {
    bookID := c.Param("id")
    userID := c.GetString("user_id")

    var req dto.CreateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    comment, err := h.bookService.AddComment(userID, bookID, req.Content, req.ParentID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, comment)
}

// GetComments godoc
// @Summary Get book comments
// @Description Get all comments for a book
// @Tags comments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {array} dto.CommentResponse
// @Failure 500 {object} map[string]string
// @Router /books/{id}/comments [get]
func (h *BookHandler) GetComments(c *gin.Context) {
    bookID := c.Param("id")

    comments, err := h.bookService.GetComments(bookID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, comments)
}

// CreateCategory godoc
// @Summary Create a category (Owner only)
// @Description Create a new category (Owner only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Category Request"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]string
// @Router /categories [post]
func (h *BookHandler) CreateCategory(c *gin.Context) {
    var req dto.CreateCategoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    category, err := h.bookService.CreateCategory(&req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, category)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Get all available categories
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CategoryResponse
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func (h *BookHandler) GetAllCategories(c *gin.Context) {
    categories, err := h.bookService.GetAllCategories()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, categories)
}