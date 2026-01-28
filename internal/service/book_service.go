package service

import (
    "errors"
    "library-project/internal/models"
    "library-project/internal/dto"
    "library-project/internal/repository"
)

type BookService struct {
    bookRepo     *repository.BookRepository
    categoryRepo *repository.CategoryRepository
    likeRepo     *repository.LikeRepository
    savedRepo    *repository.SavedBookRepository
    commentRepo  *repository.CommentRepository
}

func NewBookService(
    bookRepo *repository.BookRepository,
    categoryRepo *repository.CategoryRepository,
    likeRepo *repository.LikeRepository,
    savedRepo *repository.SavedBookRepository,
    commentRepo *repository.CommentRepository,
) *BookService {
    return &BookService{
        bookRepo:     bookRepo,
        categoryRepo: categoryRepo,
        likeRepo:     likeRepo,
        savedRepo:    savedRepo,
        commentRepo:  commentRepo,
    }
}

func (s *BookService) CreateBook(req *dto.CreateBookRequest, pdfFile, ownerID string) (*models.Book, error) {
    category, err := s.categoryRepo.FindByID(req.CategoryID)
    if err != nil {
        return nil, err
    }
    if category == nil {
        return nil, errors.New("category not found")
    }

    book := &models.Book{
        Title:       req.Title,
        Description: req.Description,
        PDFFile:     pdfFile,
        CategoryID:  req.CategoryID,
        OwnerID:     ownerID,
    }

    if err := s.bookRepo.Create(book); err != nil {
        return nil, err
    }

    return book, nil
}

func (s *BookService) UpdateBook(id string, req *dto.UpdateBookRequest, ownerID string) error {
    book, err := s.bookRepo.FindByID(id)
    if err != nil {
        return err
    }
    if book == nil {
        return errors.New("book not found")
    }
    if book.OwnerID != ownerID {
        return errors.New("unauthorized")
    }

    if req.CategoryID != "" {
        category, err := s.categoryRepo.FindByID(req.CategoryID)
        if err != nil {
            return err
        }
        if category == nil {
            return errors.New("category not found")
        }
    }

    updatedBook := &models.Book{
        ID:          id,
        Title:       req.Title,
        Description: req.Description,
        CategoryID:  req.CategoryID,
    }

    return s.bookRepo.Update(updatedBook)
}

func (s *BookService) DeleteBook(id, ownerID string) error {
    book, err := s.bookRepo.FindByID(id)
    if err != nil {
        return err
    }
    if book == nil {
        return errors.New("book not found")
    }
    if book.OwnerID != ownerID {
        return errors.New("unauthorized")
    }

    return s.bookRepo.Delete(id)
}

func (s *BookService) GetBook(id string) (*models.BookWithCategory, error) {
    return s.bookRepo.FindByID(id)
}

func (s *BookService) GetAllBooks() ([]*models.BookWithCategory, error) {
    return s.bookRepo.FindAll()
}

func (s *BookService) GetBooksByCategory(categoryID string) ([]*models.BookWithCategory, error) {
    return s.bookRepo.FindByCategory(categoryID)
}

func (s *BookService) SaveBook(userID, bookID string) error {
    book, err := s.bookRepo.FindByID(bookID)
    if err != nil {
        return err
    }
    if book == nil {
        return errors.New("book not found")
    }

    saved := &models.SavedBook{
        UserID: userID,
        BookID: bookID,
    }

    if err := s.savedRepo.Create(saved); err != nil {
        return err
    }

    return s.bookRepo.UpdateSaveCount(bookID)
}

func (s *BookService) UnsaveBook(userID, bookID string) error {
    if err := s.savedRepo.Delete(userID, bookID); err != nil {
        return err
    }

    return s.bookRepo.UpdateSaveCount(bookID)
}

func (s *BookService) GetSavedBooks(userID string) ([]*models.BookWithCategory, error) {
    return s.savedRepo.FindByUserID(userID)
}

func (s *BookService) LikeBook(userID, bookID string, isLike bool) error {
    book, err := s.bookRepo.FindByID(bookID)
    if err != nil {
        return err
    }
    if book == nil {
        return errors.New("book not found")
    }

    like := &models.Like{
        UserID: userID,
        BookID: bookID,
        IsLike: isLike,
    }

    if err := s.likeRepo.Upsert(like); err != nil {
        return err
    }

    return s.bookRepo.UpdateLikeCount(bookID)
}

func (s *BookService) RemoveLike(userID, bookID string) error {
    if err := s.likeRepo.Delete(userID, bookID); err != nil {
        return err
    }

    return s.bookRepo.UpdateLikeCount(bookID)
}

func (s *BookService) AddComment(userID, bookID, content string, parentID *string) (*models.Comment, error) {
    book, err := s.bookRepo.FindByID(bookID)
    if err != nil {
        return nil, err
    }
    if book == nil {
        return nil, errors.New("book not found")
    }

    comment := &models.Comment{
        BookID:   bookID,
        UserID:   userID,
        Content:  content,
        ParentID: parentID,
    }

    if err := s.commentRepo.Create(comment); err != nil {
        return nil, err
    }

    return comment, nil
}

func (s *BookService) GetComments(bookID string) ([]*models.CommentWithUser, error) {
    return s.commentRepo.FindByBookID(bookID)
}

func (s *BookService) CreateCategory(req *dto.CreateCategoryRequest) (*models.Category, error) {
    category := &models.Category{
        Name:        req.Name,
        Description: req.Description,
    }

    if err := s.categoryRepo.Create(category); err != nil {
        return nil, err
    }

    return category, nil
}

func (s *BookService) GetAllCategories() ([]*models.Category, error) {
    return s.categoryRepo.FindAll()
}