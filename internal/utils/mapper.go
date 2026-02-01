package utils

import (
	"library-project/internal/dto"
	"library-project/internal/models"
)

// MapUserToResponse converts User model to UserResponse DTO
func MapUserToResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

// // MapBookToResponse converts Book model to BookResponse DTO
// func MapBookToResponse(book *models.BookWithCategory) dto.BookResponse {
// 	return dto.BookResponse{
// 		ID:           book.ID,
// 		Title:        book.Title,
// 		Description:  book.Description,
// 		PDFFile:      book.PDFFile,
// 		CategoryID:   book.CategoryID,
// 		CategoryName: book.CategoryName,
// 		OwnerID:      book.OwnerID,
// 		LikeCount:    book.LikeCount,
// 		DislikeCount: book.DislikeCount,
// 		SaveCount:    book.SaveCount,
// 		CreatedAt:    book.CreatedAt,
// 		UpdatedAt:    book.UpdatedAt,
// 	}
// }

// // MapBooksToResponse converts slice of Book models to BookResponse DTOs
// func MapBooksToResponse(books []*models.BookWithCategory) []dto.BookResponse {
// 	responses := make([]dto.BookResponse, len(books))
// 	for i, book := range books {
// 		responses[i] = MapBookToResponse(book)
// 	}
// 	return responses
// }

// // MapCategoryToResponse converts Category model to CategoryResponse DTO
// func MapCategoryToResponse(category *models.Category) dto.CategoryResponse {
// 	return dto.CategoryResponse{
// 		ID:          category.ID,
// 		Name:        category.Name,
// 		Description: category.Description,
// 		CreatedAt:   category.CreatedAt,
// 		UpdatedAt:   category.UpdatedAt,
// 	}
// }

// // MapCategoriesToResponse converts slice of Category models to CategoryResponse DTOs
// func MapCategoriesToResponse(categories []*models.Category) []dto.CategoryResponse {
// 	responses := make([]dto.CategoryResponse, len(categories))
// 	for i, cat := range categories {
// 		responses[i] = MapCategoryToResponse(cat)
// 	}
// 	return responses
// }

// // MapCommentToResponse converts Comment model to CommentResponse DTO
// func MapCommentToResponse(comment *models.CommentWithUser) dto.CommentResponse {
// 	return dto.CommentResponse{
// 		ID:            comment.ID,
// 		BookID:        comment.BookID,
// 		UserID:        comment.UserID,
// 		UserFirstName: comment.UserFirstName,
// 		UserLastName:  comment.UserLastName,
// 		Content:       comment.Content,
// 		ParentID:      comment.ParentID,
// 		CreatedAt:     comment.CreatedAt,
// 		UpdatedAt:     comment.UpdatedAt,
// 	}
// }

// // MapCommentsToResponse converts slice of Comment models to CommentResponse DTOs
// func MapCommentsToResponse(comments []*models.CommentWithUser) []dto.CommentResponse {
// 	responses := make([]dto.CommentResponse, len(comments))
// 	for i, comment := range comments {
// 		responses[i] = MapCommentToResponse(comment)
// 	}
// 	return responses
// }