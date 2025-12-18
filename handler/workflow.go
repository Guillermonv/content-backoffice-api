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

func RegisterWorkflowRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config) {
    g := r.Group("/workflows")
    // Aplicar autenticación JWT a todas las rutas
    g.Use(middleware.AuthMiddleware(cfg))

    // Rutas de lectura requieren scope workflows:read
    g.GET("", middleware.RequireScopes("workflows:read"), func(c *gin.Context) {
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

        var list []model.Workflow
        db.Scopes(service.Paginate(page, size)).Find(&list)
        c.JSON(http.StatusOK, list)
    })

    // Rutas de escritura requieren scope workflows:write
    g.POST("", middleware.RequireScopes("workflows:write"), func(c *gin.Context) {
        var w model.Workflow
        if err := c.ShouldBindJSON(&w); err != nil {
            c.JSON(http.StatusBadRequest, err)
            return
        }
        db.Create(&w)
        c.JSON(http.StatusCreated, w)
    })

    g.PUT("/:id", middleware.RequireScopes("workflows:write"), func(c *gin.Context) {
        id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

        var w model.Workflow
        if err := db.First(&w, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
            return
        }

        var input model.Workflow
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, err)
            return
        }

        db.Model(&w).Updates(input)
        c.JSON(http.StatusOK, w)
    })

    g.DELETE("/:id", middleware.RequireScopes("workflows:write"), func(c *gin.Context) {
        id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
        db.Delete(&model.Workflow{}, id)
        c.Status(http.StatusNoContent)
    })
}
