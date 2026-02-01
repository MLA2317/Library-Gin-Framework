package middleware

import (
    "io"
    "library-project/config"
    "library-project/internal/models"
    "library-project/internal/repository"
    "library-project/internal/utils"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

// drainBody reads and discards the request body to prevent connection issues
// when aborting multipart/form-data requests early.
func drainBody(c *gin.Context) {
    if c.Request.Body != nil {
        io.Copy(io.Discard, c.Request.Body)
        c.Request.Body.Close()
    }
}

func AuthMiddleware(cfg *config.Config, userRepo ...*repository.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            drainBody(c)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            drainBody(c)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
            c.Abort()
            return
        }

        claims, err := utils.ValidateToken(parts[1], cfg.JWT.Secret)
        if err != nil {
            drainBody(c)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
            c.Abort()
            return
        }

        // Check if token was invalidated (logout/refresh)
        if len(userRepo) > 0 && userRepo[0] != nil {
            user, err := userRepo[0].FindByID(claims.UserID)
            if err != nil || user == nil {
                drainBody(c)
                c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
                c.Abort()
                return
            }

            if user.TokenInvalidatedAt != nil && claims.IssuedAt != nil {
                if claims.IssuedAt.Time.Before(*user.TokenInvalidatedAt) {
                    drainBody(c)
                    c.JSON(http.StatusUnauthorized, gin.H{"error": "token has been invalidated"})
                    c.Abort()
                    return
                }
            }
        }

        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("user_role", claims.Role)
        c.Next()
    }
}

func RoleMiddleware(allowedRoles ...models.UserRole) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("user_role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "user role not found"})
            c.Abort()
            return
        }

        userRole := models.UserRole(role.(string))
        for _, allowedRole := range allowedRoles {
            if userRole == allowedRole {
                c.Next()
                return
            }
        }

        drainBody(c)
        c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
        c.Abort()
    }
}
