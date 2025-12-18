package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"example.com/workflowapi/config"
	"example.com/workflowapi/middleware"
	"example.com/workflowapi/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password" binding:"required"`
	Scopes   []string `json:"scopes,omitempty"` // Scopes opcionales en la request
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id"`
	Scopes    []string  `json:"scopes"`
}

// RegisterAuthRoutes registra las rutas de autenticación
func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", loginHandler(db, cfg))
		auth.POST("/validate", middleware.AuthMiddleware(cfg), validateTokenHandler())
	}
}

// loginHandler genera un token JWT después de validar credenciales contra la base de datos
func loginHandler(db *gorm.DB, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// VALIDACIÓN ESTRICTA: Primero verificar que el usuario existe
		var count int64
		if err := db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
			log.Printf("Database error checking if user '%s' exists: %v", req.Username, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Si no existe ningún usuario con ese nombre, rechazar inmediatamente
		if count == 0 {
			log.Printf("Login rejected: user '%s' does not exist in database", req.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Ahora buscar el usuario para obtener sus datos
		var user model.User
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			// Esto no debería pasar ya que verificamos que existe, pero por seguridad lo manejamos
			log.Printf("Unexpected error fetching user '%s': %v", req.Username, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Verificación adicional: el usuario DEBE tener datos válidos
		if user.ID == 0 || user.Username == "" || user.PasswordHash == "" {
			log.Printf("Login rejected: user '%s' found but has invalid data (ID=%d)", req.Username, user.ID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		log.Printf("User '%s' (ID: %d) found in database - validating password", user.Username, user.ID)

		// Verificar que el usuario esté activo
		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is disabled"})
			return
		}

		// Verificar contraseña - SIEMPRE validar contra la BD
		if !user.CheckPassword(req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Obtener scopes del usuario DESDE LA BASE DE DATOS (nunca de la request por seguridad)
		var scopes []string
		if user.Scopes != "" {
			if err := json.Unmarshal([]byte(user.Scopes), &scopes); err != nil {
				// Si hay error parseando scopes, usar scopes vacíos (sin permisos)
				scopes = []string{}
			}
		}

		// Si el usuario no tiene scopes asignados, NO permitir acceso (sin permisos)
		if len(scopes) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "User has no permissions assigned"})
			return
		}

		// Crear claims del token
		expirationTime := time.Now().Add(24 * time.Hour) // Token válido por 24 horas
		claims := &middleware.Claims{
			UserID: req.Username,
			Scopes: scopes,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "n-backoffice-api",
				Subject:   req.Username,
			},
		}

		// Generar token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token:     tokenString,
			ExpiresAt: expirationTime,
			UserID:    req.Username,
			Scopes:    scopes,
		})
	}
}

// validateTokenHandler valida que un token sea válido
func validateTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		scopes, _ := c.Get("scopes")

		c.JSON(http.StatusOK, gin.H{
			"valid":   true,
			"user_id": userID,
			"scopes":  scopes,
		})
	}
}
