package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"github.com/hadi-projects/go-react-starter/internal/generator"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/response"
	"gorm.io/gorm"
)

type GeneratorHandler interface {
	Generate(c *gin.Context)
}

type generatorHandler struct {
	baseDir string
	db      *gorm.DB
}

func NewGeneratorHandler(baseDir string, db *gorm.DB) GeneratorHandler {
	return &generatorHandler{
		baseDir: baseDir,
		db:      db,
	}
}

func (h *generatorHandler) Generate(c *gin.Context) {
	var config generator.ModuleConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	g := generator.NewGeneratorFromConfig(config, h.baseDir)
	if err := g.Generate(); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Failed to generate module")
		response.Error(c, http.StatusInternalServerError, "Failed to generate module: "+err.Error())
		return
	}
	// Generate basic permissions and assign to admin role
	if err := h.generatePermissions(config.ModuleName); err != nil {
		logger.SystemLogger.Error().Err(err).Msg("Failed to auto-generate permissions for module")
		// Log error but don't fail the generator response
	}

	response.Success(c, http.StatusOK, "Module generated successfully", nil)
}

func (h *generatorHandler) generatePermissions(moduleName string) error {
	if h.db == nil {
		return fmt.Errorf("db not initialized in generator handler")
	}

	baseName := strings.ToLower(moduleName)
	permissions := []entity.Permission{
		{Name: "get-" + baseName, Description: "Read access for " + moduleName},
		{Name: "create-" + baseName, Description: "Create access for " + moduleName},
		{Name: "update-" + baseName, Description: "Update access for " + moduleName},
		{Name: "delete-" + baseName, Description: "Delete access for " + moduleName},
	}

	// Insert permissions
	for i := range permissions {
		if err := h.db.FirstOrCreate(&permissions[i], entity.Permission{Name: permissions[i].Name}).Error; err != nil {
			return err
		}
	}

	// Assign to admin role (Role ID = 1)
	var adminRole entity.Role
	if err := h.db.Preload("Permissions").First(&adminRole, 1).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // No admin role, safe to ignore
		}
		return err
	}

	// Add new permissions to role if not already there
	for _, p := range permissions {
		exists := false
		for _, rp := range adminRole.Permissions {
			if rp.ID == p.ID {
				exists = true
				break
			}
		}
		if !exists {
			adminRole.Permissions = append(adminRole.Permissions, p)
		}
	}

	return h.db.Save(&adminRole).Error
}
