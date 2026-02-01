package utils

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// ValidatePDFFile checks if the uploaded file is a valid PDF
func ValidatePDFFile(file *multipart.FileHeader, maxSize int64) error {
	if file == nil {
		return errors.New("file is required")
	}

	// Check file size
	if file.Size > maxSize {
		return errors.New("file size exceeds maximum limit")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		return errors.New("only PDF files are allowed")
	}

	// Check MIME type by reading file header
	src, err := file.Open()
	if err != nil {
		return errors.New("failed to open file")
	}
	defer src.Close()

	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return errors.New("failed to read file")
	}

	// Detect content type
	contentType := http.DetectContentType(buffer)
	if contentType != "application/pdf" {
		return errors.New("invalid file type: file must be a PDF")
	}

	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 72 {
		return errors.New("password must not exceed 72 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	// Simple email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}
