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

func RegisterNRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config) {
    g := r.Group("/n")
    // Aplicar autenticación JWT a todas las rutas
    g.Use(middleware.AuthMiddleware(cfg))

    // Rutas de lectura requieren scope n:read
    g.GET("", middleware.RequireScopes("n:read"), func(c *gin.Context) {
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

        var list []model.N
        db.Scopes(service.Paginate(page, size)).Find(&list)
        c.JSON(http.StatusOK, list)
    })

    // Rutas de escritura requieren scope n:write
    g.POST("", middleware.RequireScopes("n:write"), func(c *gin.Context) {
        var n model.N
        if err := c.ShouldBindJSON(&n); err != nil {
            c.JSON(http.StatusBadRequest, err)
            return
        }
        db.Create(&n)
        c.JSON(http.StatusCreated, n)
    })

    g.PUT("/:id", middleware.RequireScopes("n:write"), func(c *gin.Context) {
        id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

        var n model.N
        if err := db.First(&n, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
            return
        }

        var input model.N
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, err)
            return
        }

        db.Model(&n).Updates(input)
        c.JSON(http.StatusOK, n)
    })

    g.DELETE("/:id", middleware.RequireScopes("n:write"), func(c *gin.Context) {
        id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
        db.Delete(&model.N{}, id)
        c.Status(http.StatusNoContent)
    })
}
