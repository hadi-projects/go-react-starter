package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/internal/generator"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type GeneratorHandler interface {
	Generate(c *gin.Context)
}

type generatorHandler struct {
	baseDir string
}

func NewGeneratorHandler(baseDir string) GeneratorHandler {
	return &generatorHandler{
		baseDir: baseDir,
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

	response.Success(c, http.StatusOK, "Module generated successfully", nil)
}
