package utils

import (
	"golang.org/x/crypto/bcrypt"
	"library-project/internal/dto"
	"math"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// BuildPaginationResponse creates a pagination response
func BuildPaginationResponse(page, pageSize, total int) dto.PaginationResponse {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return dto.PaginationResponse{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		TotalItems:  int64(total),
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}
}