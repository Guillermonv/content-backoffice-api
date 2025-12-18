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

    g.PUT("/:id", middleware.RequireScopes("agents:write"), func(c *gin.Context) {
        id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

        var a model.Agent
        if err := db.First(&a, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
            return
        }

        var input model.Agent
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, err)
            return
        }

        db.Model(&a).Updates(input)
        c.JSON(http.StatusOK, a)
    })

    g.DELETE("/:id", middleware.RequireScopes("agents:write"), func(c *gin.Context) {
        id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
        db.Delete(&model.Agent{}, id)
        c.Status(http.StatusNoContent)
    })
}
