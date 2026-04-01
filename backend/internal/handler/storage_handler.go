package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/service"
	"github.com/hadi-projects/go-react-starter/pkg/response"
)

type StorageHandler interface {
	// Authenticated endpoints
	Upload(c *gin.Context)
	GetFiles(c *gin.Context)
	GetFileByID(c *gin.Context)
	DeleteFile(c *gin.Context)
	DownloadFile(c *gin.Context)

	// Share link management (authenticated)
	CreateShareLink(c *gin.Context)
	GetShareLinks(c *gin.Context)
	UpdateShareLink(c *gin.Context)
	RevokeShareLink(c *gin.Context)
	GetShareLinkLogs(c *gin.Context)

	// Public endpoints (no auth)
	PublicFileInfo(c *gin.Context)
	PublicDownload(c *gin.Context)
}

type storageHandler struct {
	service service.StorageService
}

func NewStorageHandler(svc service.StorageService) StorageHandler {
	return &storageHandler{service: svc}
}

// ─── Authenticated ────────────────────────────────────────────────────────────

func (h *storageHandler) Upload(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file is required")
		return
	}

	description := c.PostForm("description")

	res, err := h.service.Upload(c.Request.Context(), uid, fileHeader, description)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "File uploaded successfully", res)
}

func (h *storageHandler) GetFiles(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var pagination defaultDto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.GetFiles(c.Request.Context(), uid, &pagination)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Files retrieved successfully", res)
}

func (h *storageHandler) GetFileByID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	id, _ := strconv.Atoi(c.Param("id"))

	res, err := h.service.GetFileByID(c.Request.Context(), uint(id), uid)
	if err != nil {
		response.Error(c, http.StatusNotFound, "file not found")
		return
	}
	response.Success(c, http.StatusOK, "File retrieved successfully", res)
}

func (h *storageHandler) DeleteFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.DeleteFile(c.Request.Context(), uint(id), uid); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "file not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "File deleted successfully", nil)
}

func (h *storageHandler) DownloadFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	id, _ := strconv.Atoi(c.Param("id"))

	reader, file, err := h.service.GetFileForDownload(c.Request.Context(), uint(id), uid)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "file not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.OriginalName))
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatInt(file.Size, 10))
	c.Header("Cache-Control", "no-store")
	c.DataFromReader(http.StatusOK, file.Size, file.MimeType, reader, nil)
}

// ─── Share link management ────────────────────────────────────────────────────

func (h *storageHandler) CreateShareLink(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	fileID, _ := strconv.Atoi(c.Param("id"))

	var req dto.CreateShareLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.CreateShareLink(c.Request.Context(), uint(fileID), uid, req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "file not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Share link created", res)
}

func (h *storageHandler) GetShareLinks(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	fileID, _ := strconv.Atoi(c.Param("id"))

	res, err := h.service.GetShareLinks(c.Request.Context(), uint(fileID), uid)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "file not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Share links retrieved", res)
}

func (h *storageHandler) UpdateShareLink(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	shareID, _ := strconv.Atoi(c.Param("shareId"))

	var req dto.UpdateShareLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateShareLink(c.Request.Context(), uint(shareID), uid, req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "share link not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Share link updated", res)
}

func (h *storageHandler) RevokeShareLink(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	shareID, _ := strconv.Atoi(c.Param("shareId"))

	if err := h.service.RevokeShareLink(c.Request.Context(), uint(shareID), uid); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "share link not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Share link revoked", nil)
}

func (h *storageHandler) GetShareLinkLogs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	shareID, _ := strconv.Atoi(c.Param("shareId"))

	res, err := h.service.GetShareLinkLogs(c.Request.Context(), uint(shareID), uid)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "share link not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Access logs retrieved", res)
}

// ─── Public endpoints ─────────────────────────────────────────────────────────

func (h *storageHandler) PublicFileInfo(c *gin.Context) {
	token := c.Param("token")

	res, err := h.service.GetPublicFileInfo(c.Request.Context(), token)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "link not found or expired")
			return
		}
		response.Error(c, http.StatusForbidden, "this link is no longer accessible")
		return
	}
	response.Success(c, http.StatusOK, "File info retrieved", res)
}

func (h *storageHandler) PublicDownload(c *gin.Context) {
	token := c.Param("token")
	password := c.Query("password")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	reader, file, allowDownload, err := h.service.ServePublicFile(
		c.Request.Context(), token, password, ip, userAgent,
	)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "link not found or expired")
			return
		}
		if errors.Is(err, service.ErrInvalidPassword) {
			response.Error(c, http.StatusUnauthorized, "incorrect password")
			return
		}
		response.Error(c, http.StatusForbidden, "this link is no longer accessible")
		return
	}
	defer reader.Close()

	disposition := "inline"
	if allowDownload {
		disposition = fmt.Sprintf(`attachment; filename="%s"`, file.OriginalName)
	} else {
		disposition = fmt.Sprintf(`inline; filename="%s"`, file.OriginalName)
	}

	c.Header("Content-Disposition", disposition)
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatInt(file.Size, 10))
	c.Header("Cache-Control", "no-store")
	c.DataFromReader(http.StatusOK, file.Size, file.MimeType, reader, nil)
}
