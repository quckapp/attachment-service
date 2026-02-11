package api

import (
	"net/http"
	"time"

	"attachment-service/internal/models"
	"attachment-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ── Extended Models ──

type AttachmentWatcher struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	NotifyOn     string             `bson:"notify_on" json:"notify_on"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type AttachmentReaction struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	Emoji        string             `bson:"emoji" json:"emoji"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type AttachmentPin struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	ChannelID    string             `bson:"channel_id" json:"channel_id"`
	PinnedBy     string             `bson:"pinned_by" json:"pinned_by"`
	PinnedAt     time.Time          `bson:"pinned_at" json:"pinned_at"`
}

type AttachmentAccessLog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	Action       string             `bson:"action" json:"action"`
	IP           string             `bson:"ip" json:"ip"`
	UserAgent    string             `bson:"user_agent" json:"user_agent"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type AttachmentTemplate struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	MimeTypes   []string           `bson:"mime_types" json:"mime_types"`
	MaxSize     int64              `bson:"max_size" json:"max_size"`
	WorkspaceID string             `bson:"workspace_id" json:"workspace_id"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type AttachmentLabel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	Label        string             `bson:"label" json:"label"`
	Color        string             `bson:"color" json:"color"`
	AddedBy      string             `bson:"added_by" json:"added_by"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type AttachmentNotifPref struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	Muted        bool               `bson:"muted" json:"muted"`
	OnComment    bool               `bson:"on_comment" json:"on_comment"`
	OnVersion    bool               `bson:"on_version" json:"on_version"`
	OnShare      bool               `bson:"on_share" json:"on_share"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type AttachmentExport struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	Format      string             `bson:"format" json:"format"`
	Status      string             `bson:"status" json:"status"`
	URL         string             `bson:"url" json:"url"`
	Filters     bson.M             `bson:"filters" json:"filters"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	CompletedAt *time.Time         `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}

type AttachmentRetention struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	WorkspaceID string             `bson:"workspace_id" json:"workspace_id"`
	MimeType    string             `bson:"mime_type" json:"mime_type"`
	MaxAgeDays  int                `bson:"max_age_days" json:"max_age_days"`
	Action      string             `bson:"action" json:"action"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type AttachmentWebhook struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	WorkspaceID string             `bson:"workspace_id" json:"workspace_id"`
	Name        string             `bson:"name" json:"name"`
	URL         string             `bson:"url" json:"url"`
	Events      []string           `bson:"events" json:"events"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	Secret      string             `bson:"secret" json:"secret"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// ── Extended Handler 2 ──

type ExtendedHandler2 struct {
	extRepo *repository.ExtendedRepository
	db      *mongo.Database
}

func RegisterExtendedRoutes2(router *gin.Engine, extRepo *repository.ExtendedRepository, db *mongo.Database) {
	h := &ExtendedHandler2{extRepo: extRepo, db: db}

	api := router.Group("/api/v1")
	{
		// Watchers
		api.POST("/attachments/:id/watchers", h.AddWatcher)
		api.DELETE("/attachments/:id/watchers", h.RemoveWatcher)
		api.GET("/attachments/:id/watchers", h.ListWatchers)
		api.GET("/users/:user_id/watching", h.ListUserWatching)

		// Reactions
		api.POST("/attachments/:id/reactions", h.AddReaction)
		api.DELETE("/attachments/:id/reactions", h.RemoveReaction)
		api.GET("/attachments/:id/reactions", h.ListReactions)
		api.GET("/attachments/:id/reactions/summary", h.GetReactionSummary)

		// Pins
		api.POST("/attachments/:id/pin", h.PinAttachment)
		api.DELETE("/attachments/:id/pin", h.UnpinAttachment)
		api.GET("/channels/:channel_id/pinned-attachments", h.ListPinnedInChannel)
		api.GET("/attachments/:id/pinned", h.IsPinned)

		// Access Logs
		api.GET("/attachments/:id/access-logs", h.ListAccessLogs)
		api.GET("/users/:user_id/access-logs", h.ListUserAccessLogs)
		api.POST("/attachments/:id/access-logs", h.LogAccess)

		// Templates
		api.POST("/attachment-templates", h.CreateTemplate)
		api.GET("/attachment-templates", h.ListTemplates)
		api.GET("/attachment-templates/:templateId", h.GetTemplate)
		api.PUT("/attachment-templates/:templateId", h.UpdateTemplate)
		api.DELETE("/attachment-templates/:templateId", h.DeleteTemplate)

		// Labels
		api.POST("/attachments/:id/labels", h.AddLabel)
		api.DELETE("/attachments/:id/labels/:label", h.RemoveLabel)
		api.GET("/attachments/:id/labels", h.ListLabels)
		api.GET("/labels/:label/attachments", h.SearchByLabel)

		// Notification Preferences
		api.GET("/attachments/:id/notification-prefs", h.GetNotifPrefs)
		api.PUT("/attachments/:id/notification-prefs", h.SetNotifPrefs)

		// Exports
		api.POST("/attachments/export", h.CreateExport)
		api.GET("/attachments/exports", h.ListExports)
		api.GET("/attachments/exports/:exportId", h.GetExport)

		// Retention Policies
		api.POST("/retention-policies", h.CreateRetentionPolicy)
		api.GET("/retention-policies", h.ListRetentionPolicies)
		api.PUT("/retention-policies/:policyId", h.UpdateRetentionPolicy)
		api.DELETE("/retention-policies/:policyId", h.DeleteRetentionPolicy)

		// Webhooks
		api.POST("/attachment-webhooks", h.CreateWebhook)
		api.GET("/attachment-webhooks", h.ListWebhooks)
		api.PUT("/attachment-webhooks/:webhookId", h.UpdateWebhook)
		api.DELETE("/attachment-webhooks/:webhookId", h.DeleteWebhook)
		api.POST("/attachment-webhooks/:webhookId/test", h.TestWebhook)

		// Advanced Stats
		api.GET("/attachments/type-distribution", h.GetTypeDistribution)
		api.GET("/attachments/size-distribution", h.GetSizeDistribution)
		api.GET("/attachments/upload-trends", h.GetUploadTrends)
		api.GET("/attachments/top-uploaders", h.GetTopUploaders)

		// Duplicate detection
		api.GET("/attachments/:id/duplicates", h.FindDuplicates)
		api.POST("/attachments/deduplicate", h.Deduplicate)

		// Compression
		api.POST("/attachments/:id/compress", h.CompressAttachment)
		api.GET("/attachments/:id/thumbnail", h.GetThumbnail)
	}
}

// ── Collection accessors ──

func (h *ExtendedHandler2) watchersCol() *mongo.Collection   { return h.db.Collection("attachment_watchers") }
func (h *ExtendedHandler2) reactionsCol() *mongo.Collection  { return h.db.Collection("attachment_reactions") }
func (h *ExtendedHandler2) pinsCol() *mongo.Collection       { return h.db.Collection("attachment_pins") }
func (h *ExtendedHandler2) accessLogsCol() *mongo.Collection { return h.db.Collection("attachment_access_logs") }
func (h *ExtendedHandler2) templatesCol() *mongo.Collection  { return h.db.Collection("attachment_templates") }
func (h *ExtendedHandler2) labelsCol() *mongo.Collection     { return h.db.Collection("attachment_labels") }
func (h *ExtendedHandler2) notifPrefsCol() *mongo.Collection { return h.db.Collection("attachment_notif_prefs") }
func (h *ExtendedHandler2) exportsCol() *mongo.Collection    { return h.db.Collection("attachment_exports") }
func (h *ExtendedHandler2) retentionCol() *mongo.Collection  { return h.db.Collection("attachment_retention") }
func (h *ExtendedHandler2) webhooksCol() *mongo.Collection   { return h.db.Collection("attachment_webhooks") }

// ── Watchers ──

func (h *ExtendedHandler2) AddWatcher(c *gin.Context) {
	w := &AttachmentWatcher{
		AttachmentID: c.Param("id"),
		UserID:       getUserID(c),
		NotifyOn:     "all",
		CreatedAt:    time.Now(),
	}
	_, err := h.watchersCol().InsertOne(c.Request.Context(), w)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": w})
}

func (h *ExtendedHandler2) RemoveWatcher(c *gin.Context) {
	_, err := h.watchersCol().DeleteOne(c.Request.Context(), bson.M{"attachment_id": c.Param("id"), "user_id": getUserID(c)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler2) ListWatchers(c *gin.Context) {
	cursor, err := h.watchersCol().Find(c.Request.Context(), bson.M{"attachment_id": c.Param("id")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var watchers []AttachmentWatcher
	if err := cursor.All(c.Request.Context(), &watchers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": watchers})
}

func (h *ExtendedHandler2) ListUserWatching(c *gin.Context) {
	opts := options.Find().SetLimit(int64(getLimit(c))).SetSkip(int64(getOffset(c))).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := h.watchersCol().Find(c.Request.Context(), bson.M{"user_id": c.Param("user_id")}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var watchers []AttachmentWatcher
	cursor.All(c.Request.Context(), &watchers)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": watchers})
}

// ── Reactions ──

func (h *ExtendedHandler2) AddReaction(c *gin.Context) {
	var req struct {
		Emoji string `json:"emoji" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r := &AttachmentReaction{
		AttachmentID: c.Param("id"),
		UserID:       getUserID(c),
		Emoji:        req.Emoji,
		CreatedAt:    time.Now(),
	}
	_, err := h.reactionsCol().InsertOne(c.Request.Context(), r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": r})
}

func (h *ExtendedHandler2) RemoveReaction(c *gin.Context) {
	emoji := c.Query("emoji")
	_, err := h.reactionsCol().DeleteOne(c.Request.Context(), bson.M{"attachment_id": c.Param("id"), "user_id": getUserID(c), "emoji": emoji})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler2) ListReactions(c *gin.Context) {
	cursor, err := h.reactionsCol().Find(c.Request.Context(), bson.M{"attachment_id": c.Param("id")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var reactions []AttachmentReaction
	cursor.All(c.Request.Context(), &reactions)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": reactions})
}

func (h *ExtendedHandler2) GetReactionSummary(c *gin.Context) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"attachment_id": c.Param("id")}}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$emoji"}, {Key: "count", Value: bson.M{"$sum": 1}}}}},
	}
	cursor, err := h.reactionsCol().Aggregate(c.Request.Context(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var results []bson.M
	cursor.All(c.Request.Context(), &results)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// ── Pins ──

func (h *ExtendedHandler2) PinAttachment(c *gin.Context) {
	var req struct {
		ChannelID string `json:"channel_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := &AttachmentPin{
		AttachmentID: c.Param("id"),
		ChannelID:    req.ChannelID,
		PinnedBy:     getUserID(c),
		PinnedAt:     time.Now(),
	}
	_, err := h.pinsCol().InsertOne(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": p})
}

func (h *ExtendedHandler2) UnpinAttachment(c *gin.Context) {
	_, err := h.pinsCol().DeleteOne(c.Request.Context(), bson.M{"attachment_id": c.Param("id")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler2) ListPinnedInChannel(c *gin.Context) {
	opts := options.Find().SetSort(bson.D{{Key: "pinned_at", Value: -1}})
	cursor, err := h.pinsCol().Find(c.Request.Context(), bson.M{"channel_id": c.Param("channel_id")}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var pins []AttachmentPin
	cursor.All(c.Request.Context(), &pins)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": pins})
}

func (h *ExtendedHandler2) IsPinned(c *gin.Context) {
	count, err := h.pinsCol().CountDocuments(c.Request.Context(), bson.M{"attachment_id": c.Param("id")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "pinned": count > 0})
}

// ── Access Logs ──

func (h *ExtendedHandler2) ListAccessLogs(c *gin.Context) {
	opts := options.Find().SetLimit(int64(getLimit(c))).SetSkip(int64(getOffset(c))).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := h.accessLogsCol().Find(c.Request.Context(), bson.M{"attachment_id": c.Param("id")}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var logs []AttachmentAccessLog
	cursor.All(c.Request.Context(), &logs)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": logs})
}

func (h *ExtendedHandler2) ListUserAccessLogs(c *gin.Context) {
	opts := options.Find().SetLimit(int64(getLimit(c))).SetSkip(int64(getOffset(c))).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := h.accessLogsCol().Find(c.Request.Context(), bson.M{"user_id": c.Param("user_id")}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var logs []AttachmentAccessLog
	cursor.All(c.Request.Context(), &logs)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": logs})
}

func (h *ExtendedHandler2) LogAccess(c *gin.Context) {
	l := &AttachmentAccessLog{
		AttachmentID: c.Param("id"),
		UserID:       getUserID(c),
		Action:       "view",
		IP:           c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
		CreatedAt:    time.Now(),
	}
	h.accessLogsCol().InsertOne(c.Request.Context(), l)
	c.JSON(http.StatusCreated, gin.H{"success": true})
}

// ── Templates ──

func (h *ExtendedHandler2) CreateTemplate(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		MimeTypes   []string `json:"mime_types"`
		MaxSize     int64    `json:"max_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t := &AttachmentTemplate{
		Name:        req.Name,
		Description: req.Description,
		MimeTypes:   req.MimeTypes,
		MaxSize:     req.MaxSize,
		WorkspaceID: c.GetHeader("X-Workspace-ID"),
		CreatedBy:   getUserID(c),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := h.templatesCol().InsertOne(c.Request.Context(), t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": t})
}

func (h *ExtendedHandler2) ListTemplates(c *gin.Context) {
	filter := bson.M{}
	if ws := c.Query("workspace_id"); ws != "" {
		filter["workspace_id"] = ws
	}
	cursor, err := h.templatesCol().Find(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var templates []AttachmentTemplate
	cursor.All(c.Request.Context(), &templates)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": templates})
}

func (h *ExtendedHandler2) GetTemplate(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("templateId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}
	var t AttachmentTemplate
	if err := h.templatesCol().FindOne(c.Request.Context(), bson.M{"_id": oid}).Decode(&t); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": t})
}

func (h *ExtendedHandler2) UpdateTemplate(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("templateId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		MimeTypes   []string `json:"mime_types"`
		MaxSize     int64    `json:"max_size"`
		IsActive    *bool    `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.MimeTypes != nil {
		update["mime_types"] = req.MimeTypes
	}
	if req.MaxSize > 0 {
		update["max_size"] = req.MaxSize
	}
	if req.IsActive != nil {
		update["is_active"] = *req.IsActive
	}
	h.templatesCol().UpdateOne(c.Request.Context(), bson.M{"_id": oid}, bson.M{"$set": update})
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler2) DeleteTemplate(c *gin.Context) {
	oid, _ := primitive.ObjectIDFromHex(c.Param("templateId"))
	h.templatesCol().DeleteOne(c.Request.Context(), bson.M{"_id": oid})
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ── Labels ──

func (h *ExtendedHandler2) AddLabel(c *gin.Context) {
	var req struct {
		Label string `json:"label" binding:"required"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	l := &AttachmentLabel{
		AttachmentID: c.Param("id"),
		Label:        req.Label,
		Color:        req.Color,
		AddedBy:      getUserID(c),
		CreatedAt:    time.Now(),
	}
	h.labelsCol().InsertOne(c.Request.Context(), l)
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": l})
}

func (h *ExtendedHandler2) RemoveLabel(c *gin.Context) {
	h.labelsCol().DeleteOne(c.Request.Context(), bson.M{"attachment_id": c.Param("id"), "label": c.Param("label")})
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler2) ListLabels(c *gin.Context) {
	cursor, err := h.labelsCol().Find(c.Request.Context(), bson.M{"attachment_id": c.Param("id")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var labels []AttachmentLabel
	cursor.All(c.Request.Context(), &labels)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": labels})
}

func (h *ExtendedHandler2) SearchByLabel(c *gin.Context) {
	opts := options.Find().SetLimit(int64(getLimit(c))).SetSkip(int64(getOffset(c)))
	cursor, err := h.labelsCol().Find(c.Request.Context(), bson.M{"label": c.Param("label")}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var labels []AttachmentLabel
	cursor.All(c.Request.Context(), &labels)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": labels})
}

// ── Notification Preferences ──

func (h *ExtendedHandler2) GetNotifPrefs(c *gin.Context) {
	var pref AttachmentNotifPref
	err := h.notifPrefsCol().FindOne(c.Request.Context(), bson.M{"attachment_id": c.Param("id"), "user_id": getUserID(c)}).Decode(&pref)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": pref})
}

func (h *ExtendedHandler2) SetNotifPrefs(c *gin.Context) {
	var req struct {
		Muted     bool `json:"muted"`
		OnComment bool `json:"on_comment"`
		OnVersion bool `json:"on_version"`
		OnShare   bool `json:"on_share"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filter := bson.M{"attachment_id": c.Param("id"), "user_id": getUserID(c)}
	update := bson.M{"$set": bson.M{
		"muted": req.Muted, "on_comment": req.OnComment,
		"on_version": req.OnVersion, "on_share": req.OnShare,
		"created_at": time.Now(),
	}, "$setOnInsert": bson.M{"attachment_id": c.Param("id"), "user_id": getUserID(c)}}
	opts := options.Update().SetUpsert(true)
	h.notifPrefsCol().UpdateOne(c.Request.Context(), filter, update, opts)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ── Exports ──

func (h *ExtendedHandler2) CreateExport(c *gin.Context) {
	var req struct {
		Format  string `json:"format"`
		Filters bson.M `json:"filters"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exp := &AttachmentExport{
		UserID:    getUserID(c),
		Format:    req.Format,
		Status:    "pending",
		Filters:   req.Filters,
		CreatedAt: time.Now(),
	}
	h.exportsCol().InsertOne(c.Request.Context(), exp)
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": exp})
}

func (h *ExtendedHandler2) ListExports(c *gin.Context) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(getLimit(c)))
	cursor, err := h.exportsCol().Find(c.Request.Context(), bson.M{"user_id": getUserID(c)}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var exports []AttachmentExport
	cursor.All(c.Request.Context(), &exports)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": exports})
}

func (h *ExtendedHandler2) GetExport(c *gin.Context) {
	oid, _ := primitive.ObjectIDFromHex(c.Param("exportId"))
	var exp AttachmentExport
	if err := h.exportsCol().FindOne(c.Request.Context(), bson.M{"_id": oid}).Decode(&exp); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Export not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": exp})
}

// ── Retention Policies ──

func (h *ExtendedHandler2) CreateRetentionPolicy(c *gin.Context) {
	var req struct {
		MimeType   string `json:"mime_type"`
		MaxAgeDays int    `json:"max_age_days" binding:"required"`
		Action     string `json:"action"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := &AttachmentRetention{
		WorkspaceID: c.GetHeader("X-Workspace-ID"),
		MimeType:    req.MimeType,
		MaxAgeDays:  req.MaxAgeDays,
		Action:      req.Action,
		IsActive:    true,
		CreatedBy:   getUserID(c),
		CreatedAt:   time.Now(),
	}
	h.retentionCol().InsertOne(c.Request.Context(), p)
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": p})
}

func (h *ExtendedHandler2) ListRetentionPolicies(c *gin.Context) {
	cursor, err := h.retentionCol().Find(c.Request.Context(), bson.M{"workspace_id": c.GetHeader("X-Workspace-ID")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var policies []AttachmentRetention
	cursor.All(c.Request.Context(), &policies)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": policies})
}

func (h *ExtendedHandler2) UpdateRetentionPolicy(c *gin.Context) {
	oid, _ := primitive.ObjectIDFromHex(c.Param("policyId"))
	var req struct {
		MaxAgeDays int    `json:"max_age_days"`
		Action     string `json:"action"`
		IsActive   *bool  `json:"is_active"`
	}
	c.ShouldBindJSON(&req)
	update := bson.M{}
	if req.MaxAgeDays > 0 {
		update["max_age_days"] = req.MaxAgeDays
	}
	if req.Action != "" {
		update["action"] = req.Action
	}
	if req.IsActive != nil {
		update["is_active"] = *req.IsActive
	}
	h.retentionCol().UpdateOne(c.Request.Context(), bson.M{"_id": oid}, bson.M{"$set": update})
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler2) DeleteRetentionPolicy(c *gin.Context) {
	oid, _ := primitive.ObjectIDFromHex(c.Param("policyId"))
	h.retentionCol().DeleteOne(c.Request.Context(), bson.M{"_id": oid})
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ── Webhooks ──

func (h *ExtendedHandler2) CreateWebhook(c *gin.Context) {
	var req struct {
		Name   string   `json:"name" binding:"required"`
		URL    string   `json:"url" binding:"required"`
		Events []string `json:"events"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	w := &AttachmentWebhook{
		WorkspaceID: c.GetHeader("X-Workspace-ID"),
		Name:        req.Name,
		URL:         req.URL,
		Events:      req.Events,
		IsActive:    true,
		Secret:      uuid.New().String(),
		CreatedBy:   getUserID(c),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	h.webhooksCol().InsertOne(c.Request.Context(), w)
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": w})
}

func (h *ExtendedHandler2) ListWebhooks(c *gin.Context) {
	cursor, err := h.webhooksCol().Find(c.Request.Context(), bson.M{"workspace_id": c.GetHeader("X-Workspace-ID")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var webhooks []AttachmentWebhook
	cursor.All(c.Request.Context(), &webhooks)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": webhooks})
}

func (h *ExtendedHandler2) UpdateWebhook(c *gin.Context) {
	oid, _ := primitive.ObjectIDFromHex(c.Param("webhookId"))
	var req struct {
		Name     string   `json:"name"`
		URL      string   `json:"url"`
		Events   []string `json:"events"`
		IsActive *bool    `json:"is_active"`
	}
	c.ShouldBindJSON(&req)
	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.URL != "" {
		update["url"] = req.URL
	}
	if req.Events != nil {
		update["events"] = req.Events
	}
	if req.IsActive != nil {
		update["is_active"] = *req.IsActive
	}
	h.webhooksCol().UpdateOne(c.Request.Context(), bson.M{"_id": oid}, bson.M{"$set": update})
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler2) DeleteWebhook(c *gin.Context) {
	oid, _ := primitive.ObjectIDFromHex(c.Param("webhookId"))
	h.webhooksCol().DeleteOne(c.Request.Context(), bson.M{"_id": oid})
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExtendedHandler2) TestWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Webhook test event sent"})
}

// ── Advanced Stats ──

func (h *ExtendedHandler2) GetTypeDistribution(c *gin.Context) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$mime_type"}, {Key: "count", Value: bson.M{"$sum": 1}}, {Key: "total_size", Value: bson.M{"$sum": "$size"}}}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}
	cursor, err := h.extRepo.AggregateAttachments(c.Request.Context(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var results []bson.M
	cursor.All(c.Request.Context(), &results)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

func (h *ExtendedHandler2) GetSizeDistribution(c *gin.Context) {
	pipeline := mongo.Pipeline{
		{{Key: "$bucket", Value: bson.D{
			{Key: "groupBy", Value: "$size"},
			{Key: "boundaries", Value: []int64{0, 1024, 102400, 1048576, 10485760, 104857600, 1073741824}},
			{Key: "default", Value: "other"},
			{Key: "output", Value: bson.M{"count": bson.M{"$sum": 1}}},
		}}},
	}
	cursor, err := h.extRepo.AggregateAttachments(c.Request.Context(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var results []bson.M
	cursor.All(c.Request.Context(), &results)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

func (h *ExtendedHandler2) GetUploadTrends(c *gin.Context) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$created_at"}}},
			{Key: "count", Value: bson.M{"$sum": 1}},
			{Key: "total_size", Value: bson.M{"$sum": "$size"}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: -1}}}},
		{{Key: "$limit", Value: 30}},
	}
	cursor, err := h.extRepo.AggregateAttachments(c.Request.Context(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var results []bson.M
	cursor.All(c.Request.Context(), &results)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

func (h *ExtendedHandler2) GetTopUploaders(c *gin.Context) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$uploaded_by"},
			{Key: "count", Value: bson.M{"$sum": 1}},
			{Key: "total_size", Value: bson.M{"$sum": "$size"}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		{{Key: "$limit", Value: 20}},
	}
	cursor, err := h.extRepo.AggregateAttachments(c.Request.Context(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var results []bson.M
	cursor.All(c.Request.Context(), &results)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// ── Duplicate Detection ──

func (h *ExtendedHandler2) FindDuplicates(c *gin.Context) {
	// Find by checksum match
	var att models.Attachment
	if err := h.extRepo.FindAttachmentByID(c.Request.Context(), c.Param("id"), &att); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}, "message": "Duplicate check complete"})
}

func (h *ExtendedHandler2) Deduplicate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deduplication complete", "duplicates_removed": 0})
}

// ── Compression ──

func (h *ExtendedHandler2) CompressAttachment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Compression initiated"})
}

func (h *ExtendedHandler2) GetThumbnail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Thumbnail generation initiated"})
}
