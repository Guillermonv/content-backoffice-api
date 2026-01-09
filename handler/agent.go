package handler

import (
	"net/http"
	"strconv"

	"example.com/workflowapi/config"
	"example.com/workflowapi/middleware"
	"example.com/workflowapi/model"
	"example.com/workflowapi/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAgentRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config) {
    g := r.Group("/agents")
    // Aplicar autenticación JWT a todas las rutas
    g.Use(middleware.AuthMiddleware(cfg))

    // Rutas de lectura requieren scope agents:read
    g.GET("", middleware.RequireScopes("agents:read"), func(c *gin.Context) {
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

        var list []model.Agent
        db.Scopes(service.Paginate(page, size)).Find(&list)
        c.JSON(http.StatusOK, list)
    })

    // Rutas de escritura requieren scope agents:write
    g.POST("", middleware.RequireScopes("agents:write"), func(c *gin.Context) {
        var a model.Agent
        if err := c.ShouldBindJSON(&a); err != nil {
            c.JSON(http.StatusBadRequest, err)
            return
        }
        db.Create(&a)
        c.JSON(http.StatusCreated, a)
    })

    g.DELETE("/:id", middleware.RequireScopes("agents:write"), func(c *gin.Context) {
        id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
        db.Delete(&model.Agent{}, id)
        c.Status(http.StatusNoContent)
    })
    g.PUT("/:id", middleware.RequireScopes("agents:write"), func(c *gin.Context) {
        // 1️⃣ parsear ID
        id, err := strconv.ParseUint(c.Param("id"), 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent id"})
            return
        }
    
        // 2️⃣ buscar existente
        var existing model.Agent
        if err := db.First(&existing, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
            return
        }
    
        // 3️⃣ bind body
        var input model.Agent
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, err)
            return
        }
    
        // 4️⃣ update REAL (map → zero-values OK)
        updates := map[string]interface{}{
            "provider": input.Provider,
            "secret":  input.Secret,
        }
    
        if err := db.Model(&existing).Updates(updates).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    
        // 5️⃣ reload
        if err := db.First(&existing, existing.ID).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload agent"})
            return
        }
    
        c.JSON(http.StatusOK, existing)
    })
}    