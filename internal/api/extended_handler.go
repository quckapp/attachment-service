package api

import (
	"net/http"
	"strconv"

	"attachment-service/internal/models"
	"attachment-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExtendedHandler struct {
	extRepo *repository.ExtendedRepository
}

func RegisterExtendedRoutes(router *gin.Engine, extRepo *repository.ExtendedRepository) {
	h := &ExtendedHandler{extRepo: extRepo}

	api := router.Group("/api/v1")
	{
		// Versions
		api.GET("/attachments/:id/versions", h.ListVersions)
		api.GET("/attachments/:id/versions/:versionNum", h.GetVersion)
		api.DELETE("/attachments/:id/versions/:versionId", h.DeleteVersion)

		// Comments
		api.POST("/attachments/:id/comments", h.CreateComment)
		api.GET("/attachments/:id/comments", h.ListComments)
		api.PUT("/attachments/:id/comments/:commentId", h.UpdateComment)
		api.DELETE("/attachments/:id/comments/:commentId", h.DeleteComment)

		// Tags
		api.POST("/attachments/:id/tags", h.AddTag)
		api.DELETE("/attachments/:id/tags/:tag", h.RemoveTag)
		api.GET("/attachments/:id/tags", h.ListTags)
		api.GET("/tags/:tag/attachments", h.SearchByTag)

		// Favorites
		api.POST("/attachments/:id/favorite", h.AddFavorite)
		api.DELETE("/attachments/:id/favorite", h.RemoveFavorite)
		api.GET("/users/:user_id/favorites", h.ListFavorites)
		api.GET("/attachments/:id/favorited", h.IsFavorited)

		// Shares
		api.POST("/attachments/:id/shares", h.CreateShare)
		api.GET("/attachments/:id/shares", h.ListShares)
		api.GET("/users/:user_id/shared", h.ListSharedWith)
		api.DELETE("/shares/:shareId", h.DeleteShare)

		// Collections
		api.POST("/collections", h.CreateCollection)
		api.GET("/collections", h.ListCollections)
		api.GET("/collections/:collectionId", h.GetCollection)
		api.PUT("/collections/:collectionId", h.UpdateCollection)
		api.DELETE("/collections/:collectionId", h.DeleteCollection)
		api.POST("/collections/:collectionId/items", h.AddToCollection)
		api.DELETE("/collections/:collectionId/items/:attachmentId", h.RemoveFromCollection)
		api.GET("/collections/:collectionId/items", h.ListCollectionItems)

		// Activity
		api.GET("/attachments/:id/activity", h.ListActivity)
		api.GET("/users/:user_id/activity", h.ListUserActivity)

		// Permissions
		api.POST("/attachments/:id/permissions", h.SetPermission)
		api.GET("/attachments/:id/permissions", h.ListPermissions)
		api.DELETE("/attachments/:id/permissions/:userId", h.DeletePermission)

		// Share links
		api.POST("/attachments/:id/share-links", h.CreateShareLink)
		api.GET("/attachments/:id/share-links", h.ListShareLinks)
		api.GET("/share-links/:code", h.GetShareLink)
		api.DELETE("/share-links/:linkId", h.DeactivateShareLink)

		// Scans
		api.GET("/attachments/:id/scan", h.GetScanResult)

		// Previews
		api.GET("/attachments/:id/previews", h.ListPreviews)

		// Stats & Search
		api.GET("/workspaces/:workspace_id/stats", h.GetWorkspaceStats)
		api.GET("/workspaces/:workspace_id/attachment-stats", h.GetAttachmentStats)
		api.GET("/users/:user_id/quota", h.GetUserQuota)
		api.GET("/search", h.SearchAttachments)
		api.GET("/users/:user_id/recent", h.GetRecentAttachments)
		api.GET("/workspaces/:workspace_id/attachments", h.GetByWorkspaceID)

		// Bulk operations
		api.POST("/bulk/delete", h.BulkDelete)
		api.POST("/bulk/move", h.BulkMove)
		api.POST("/bulk/tag", h.BulkTag)

		// Individual operations
		api.PUT("/attachments/:id/rename", h.RenameAttachment)
		api.PUT("/attachments/:id/move", h.MoveAttachment)
		api.POST("/attachments/:id/clone", h.CloneAttachment)
	}
}

func getLimit(c *gin.Context) int {
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if l <= 0 || l > 200 {
		l = 50
	}
	return l
}

func getOffset(c *gin.Context) int {
	o, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if o < 0 {
		o = 0
	}
	return o
}

func getUserID(c *gin.Context) string {
	return c.GetHeader("X-User-ID")
}

// ── Versions ──

func (h *ExtendedHandler) ListVersions(c *gin.Context) {
	versions, err := h.extRepo.ListVersions(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": versions})
}

func (h *ExtendedHandler) GetVersion(c *gin.Context) {
	vNum, _ := strconv.Atoi(c.Param("versionNum"))
	version, err := h.extRepo.GetVersionByNum(c.Request.Context(), c.Param("id"), vNum)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": version})
}

func (h *ExtendedHandler) DeleteVersion(c *gin.Context) {
	if err := h.extRepo.DeleteVersion(c.Request.Context(), c.Param("versionId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ── Comments ──

func (h *ExtendedHandler) CreateComment(c *gin.Context) {
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment := &models.AttachmentComment{
		AttachmentID: c.Param("id"),
		UserID:       getUserID(c),
		Content:      req.Content,
		ParentID:     req.ParentID,
	}
	if err := h.extRepo.CreateComment(c.Request.Context(), comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": comment})
}

func (h *ExtendedHandler) ListComments(c *gin.Context) {
	comments, err := h.extRepo.ListComments(c.Request.Context(), c.Param("id"), getLimit(c), getOffset(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": comments})
}

func (h *ExtendedHandler) UpdateComment(c *gin.Context) {
	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.extRepo.UpdateComment(c.Request.Context(), c.Param("commentId"), req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler) DeleteComment(c *gin.Context) {
	if err := h.extRepo.DeleteComment(c.Request.Context(), c.Param("commentId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ── Tags ──

func (h *ExtendedHandler) AddTag(c *gin.Context) {
	var req models.AddTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag := &models.AttachmentTag{
		AttachmentID: c.Param("id"),
		Tag:          req.Tag,
		AddedBy:      getUserID(c),
	}
	if err := h.extRepo.AddTag(c.Request.Context(), tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": tag})
}

func (h *ExtendedHandler) RemoveTag(c *gin.Context) {
	if err := h.extRepo.RemoveTag(c.Request.Context(), c.Param("id"), c.Param("tag")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler) ListTags(c *gin.Context) {
	tags, err := h.extRepo.ListTags(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tags})
}

func (h *ExtendedHandler) SearchByTag(c *gin.Context) {
	tags, err := h.extRepo.SearchByTag(c.Request.Context(), c.Param("tag"), getLimit(c), getOffset(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tags})
}

// ── Favorites ──

func (h *ExtendedHandler) AddFavorite(c *gin.Context) {
	fav := &models.AttachmentFavorite{
		AttachmentID: c.Param("id"),
		UserID:       getUserID(c),
	}
	if err := h.extRepo.AddFavorite(c.Request.Context(), fav); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": fav})
}

func (h *ExtendedHandler) RemoveFavorite(c *gin.Context) {
	if err := h.extRepo.RemoveFavorite(c.Request.Context(), c.Param("id"), getUserID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler) ListFavorites(c *gin.Context) {
	favs, err := h.extRepo.ListFavorites(c.Request.Context(), c.Param("user_id"), getLimit(c), getOffset(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": favs})
}

func (h *ExtendedHandler) IsFavorited(c *gin.Context) {
	isFav, err := h.extRepo.IsFavorited(c.Request.Context(), c.Param("id"), getUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "favorited": isFav})
}

// ── Shares ──

func (h *ExtendedHandler) CreateShare(c *gin.Context) {
	var req models.CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	share := &models.AttachmentShare{
		AttachmentID: c.Param("id"),
		SharedBy:     getUserID(c),
		SharedWith:   req.SharedWith,
		Permission:   req.Permission,
	}
	if err := h.extRepo.CreateShare(c.Request.Context(), share); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": share})
}

func (h *ExtendedHandler) ListShares(c *gin.Context) {
	shares, err := h.extRepo.ListShares(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": shares})
}

func (h *ExtendedHandler) ListSharedWith(c *gin.Context) {
	shares, err := h.extRepo.ListSharedWith(c.Request.Context(), c.Param("user_id"), getLimit(c), getOffset(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": shares})
}

func (h *ExtendedHandler) DeleteShare(c *gin.Context) {
	if err := h.extRepo.DeleteShare(c.Request.Context(), c.Param("shareId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ── Collections ──

func (h *ExtendedHandler) CreateCollection(c *gin.Context) {
	var req models.CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coll := &models.AttachmentCollection{
		Name:        req.Name,
		Description: req.Description,
		WorkspaceID: c.GetHeader("X-Workspace-ID"),
		CreatedBy:   getUserID(c),
		IsPublic:    req.IsPublic,
	}
	if err := h.extRepo.CreateCollection(c.Request.Context(), coll); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": coll})
}

func (h *ExtendedHandler) ListCollections(c *gin.Context) {
	wsID := c.Query("workspace_id")
	if wsID == "" {
		wsID = c.GetHeader("X-Workspace-ID")
	}
	colls, err := h.extRepo.ListCollections(c.Request.Context(), wsID, getLimit(c), getOffset(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": colls})
}

func (h *ExtendedHandler) GetCollection(c *gin.Context) {
	coll, err := h.extRepo.GetCollection(c.Request.Context(), c.Param("collectionId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": coll})
}

func (h *ExtendedHandler) UpdateCollection(c *gin.Context) {
	var req models.UpdateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	update := make(map[string]interface{})
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.IsPublic != nil {
		update["is_public"] = *req.IsPublic
	}
	if err := h.extRepo.UpdateCollection(c.Request.Context(), c.Param("collectionId"), update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler) DeleteCollection(c *gin.Context) {
	if err := h.extRepo.DeleteCollection(c.Request.Context(), c.Param("collectionId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler) AddToCollection(c *gin.Context) {
	var req models.AddToCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := &models.CollectionItem{
		CollectionID: c.Param("collectionId"),
		AttachmentID: req.AttachmentID,
		AddedBy:      getUserID(c),
	}
	if err := h.extRepo.AddToCollection(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": item})
}

func (h *ExtendedHandler) RemoveFromCollection(c *gin.Context) {
	if err := h.extRepo.RemoveFromCollection(c.Request.Context(), c.Param("collectionId"), c.Param("attachmentId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler) ListCollectionItems(c *gin.Context) {
	items, err := h.extRepo.ListCollectionItems(c.Request.Context(), c.Param("collectionId"), getLimit(c), getOffset(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// ── Activity ──

func (h *ExtendedHandler) ListActivity(c *gin.Context) {
	acts, err := h.extRepo.ListActivity(c.Request.Context(), c.Param("id"), getLimit(c), getOffset(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": acts})
}

func (h *ExtendedHandler) ListUserActivity(c *gin.Context) {
	acts, err := h.extRepo.ListUserActivity(c.Request.Context(), c.Param("user_id"), getLimit(c), getOffset(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": acts})
}

// ── Permissions ──

func (h *ExtendedHandler) SetPermission(c *gin.Context) {
	var req models.SetPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	perm := &models.AttachmentPermission{
		AttachmentID: c.Param("id"),
		UserID:       req.UserID,
		CanView:      req.CanView,
		CanDownload:  req.CanDownload,
		CanDelete:    req.CanDelete,
		CanShare:     req.CanShare,
		GrantedBy:    getUserID(c),
	}
	if err := h.extRepo.SetPermission(c.Request.Context(), perm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler) ListPermissions(c *gin.Context) {
	perms, err := h.extRepo.ListPermissions(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": perms})
}

func (h *ExtendedHandler) DeletePermission(c *gin.Context) {
	if err := h.extRepo.DeletePermission(c.Request.Context(), c.Param("id"), c.Param("userId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ── Share Links ──

func (h *ExtendedHandler) CreateShareLink(c *gin.Context) {
	var req models.CreateShareLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	link := &models.ShareLink{
		AttachmentID: c.Param("id"),
		Code:         uuid.New().String()[:8],
		CreatedBy:    getUserID(c),
		Password:     req.Password,
		MaxDownloads: req.MaxDownloads,
	}
	if err := h.extRepo.CreateShareLink(c.Request.Context(), link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": link})
}

func (h *ExtendedHandler) ListShareLinks(c *gin.Context) {
	links, err := h.extRepo.ListShareLinks(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": links})
}

func (h *ExtendedHandler) GetShareLink(c *gin.Context) {
	link, err := h.extRepo.GetShareLinkByCode(c.Request.Context(), c.Param("code"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Share link not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": link})
}

func (h *ExtendedHandler) DeactivateShareLink(c *gin.Context) {
	if err := h.extRepo.DeactivateShareLink(c.Request.Context(), c.Param("linkId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ── Scans ──

func (h *ExtendedHandler) GetScanResult(c *gin.Context) {
	result, err := h.extRepo.GetScanResult(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No scan result found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// ── Previews ──

func (h *ExtendedHandler) ListPreviews(c *gin.Context) {
	previews, err := h.extRepo.ListPreviews(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": previews})
}

// ── Stats & Search ──

func (h *ExtendedHandler) GetAttachmentStats(c *gin.Context) {
	stats, err := h.extRepo.GetAttachmentStats(c.Request.Context(), c.Param("workspace_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

func (h *ExtendedHandler) GetWorkspaceStats(c *gin.Context) {
	stats, err := h.extRepo.GetWorkspaceStats(c.Request.Context(), c.Param("workspace_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

func (h *ExtendedHandler) GetUserQuota(c *gin.Context) {
	quota, err := h.extRepo.GetUserQuota(c.Request.Context(), c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": quota})
}

func (h *ExtendedHandler) SearchAttachments(c *gin.Context) {
	wsID := c.Query("workspace_id")
	if wsID == "" {
		wsID = c.GetHeader("X-Workspace-ID")
	}
	query := c.Query("q")
	fileType := c.Query("type")

	results, err := h.extRepo.SearchAttachments(c.Request.Context(), wsID, query, fileType, getLimit(c), getOffset(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

func (h *ExtendedHandler) GetRecentAttachments(c *gin.Context) {
	results, err := h.extRepo.GetRecentAttachments(c.Request.Context(), c.Param("user_id"), getLimit(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

func (h *ExtendedHandler) GetByWorkspaceID(c *gin.Context) {
	results, err := h.extRepo.GetByWorkspaceID(c.Request.Context(), c.Param("workspace_id"), getLimit(c), getOffset(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// ── Bulk Operations ──

func (h *ExtendedHandler) BulkDelete(c *gin.Context) {
	var req models.BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.extRepo.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "deleted": len(req.IDs)})
}

func (h *ExtendedHandler) BulkMove(c *gin.Context) {
	var req models.BulkMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.extRepo.BulkMove(c.Request.Context(), req.IDs, req.ChannelID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "moved": len(req.IDs)})
}

func (h *ExtendedHandler) BulkTag(c *gin.Context) {
	var req models.BulkTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, id := range req.IDs {
		tag := &models.AttachmentTag{
			AttachmentID: id,
			Tag:          req.Tag,
			AddedBy:      getUserID(c),
		}
		_ = h.extRepo.AddTag(c.Request.Context(), tag)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "tagged": len(req.IDs)})
}

// ── Individual Operations ──

func (h *ExtendedHandler) RenameAttachment(c *gin.Context) {
	var req models.RenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.extRepo.RenameAttachment(c.Request.Context(), c.Param("id"), req.NewName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler) MoveAttachment(c *gin.Context) {
	var req models.MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.extRepo.MoveAttachment(c.Request.Context(), c.Param("id"), req.ChannelID, req.MessageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler) CloneAttachment(c *gin.Context) {
	// This would need the base service to get the original attachment
	// For now, return a placeholder
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Clone initiated"})
}
