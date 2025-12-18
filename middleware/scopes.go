package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireScopes verifica que el usuario tenga al menos uno de los scopes requeridos
func RequireScopes(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener los scopes del contexto (seteado por AuthMiddleware)
		scopesInterface, exists := c.Get("scopes")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Scopes not found in context"})
			c.Abort()
			return
		}

		scopes, ok := scopesInterface.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid scopes format"})
			c.Abort()
			return
		}

		// Crear un map para búsqueda rápida
		scopeMap := make(map[string]bool)
		for _, scope := range scopes {
			scopeMap[scope] = true
		}

		// Verificar si el usuario tiene al menos uno de los scopes requeridos
		hasRequiredScope := false
		for _, required := range requiredScopes {
			if scopeMap[required] {
				hasRequiredScope = true
				break
			}
		}

		if !hasRequiredScope {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
				"required_scopes": requiredScopes,
				"user_scopes": scopes,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

