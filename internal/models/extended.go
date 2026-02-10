package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ── Attachment Versions ──

type AttachmentVersion struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	VersionNum   int                `bson:"version_num" json:"version_num"`
	FileName     string             `bson:"file_name" json:"file_name"`
	OriginalName string             `bson:"original_name" json:"original_name"`
	MimeType     string             `bson:"mime_type" json:"mime_type"`
	Size         int64              `bson:"size" json:"size"`
	StoragePath  string             `bson:"storage_path" json:"storage_path"`
	URL          string             `bson:"url" json:"url"`
	Checksum     string             `bson:"checksum" json:"checksum"`
	UploadedBy   string             `bson:"uploaded_by" json:"uploaded_by"`
	Comment      string             `bson:"comment" json:"comment"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// ── Attachment Comments ──

type AttachmentComment struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	Content      string             `bson:"content" json:"content"`
	ParentID     string             `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	IsEdited     bool               `bson:"is_edited" json:"is_edited"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// ── Attachment Tags ──

type AttachmentTag struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	Tag          string             `bson:"tag" json:"tag"`
	AddedBy      string             `bson:"added_by" json:"added_by"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// ── Attachment Favorites ──

type AttachmentFavorite struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// ── Attachment Shares ──

type AttachmentShare struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	SharedBy     string             `bson:"shared_by" json:"shared_by"`
	SharedWith   string             `bson:"shared_with" json:"shared_with"`
	Permission   string             `bson:"permission" json:"permission"` // view, download, edit
	ExpiresAt    *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// ── Attachment Collections / Albums ──

type AttachmentCollection struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	WorkspaceID string             `bson:"workspace_id" json:"workspace_id"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CoverURL    string             `bson:"cover_url" json:"cover_url"`
	ItemCount   int                `bson:"item_count" json:"item_count"`
	IsPublic    bool               `bson:"is_public" json:"is_public"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CollectionItem struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CollectionID string             `bson:"collection_id" json:"collection_id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	AddedBy      string             `bson:"added_by" json:"added_by"`
	Position     int                `bson:"position" json:"position"`
	AddedAt      time.Time          `bson:"added_at" json:"added_at"`
}

// ── Attachment Activity ──

type AttachmentActivity struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	Action       string             `bson:"action" json:"action"` // uploaded, downloaded, viewed, shared, deleted, commented, tagged
	Details      string             `bson:"details" json:"details"`
	IP           string             `bson:"ip" json:"ip"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// ── Attachment Permissions ──

type AttachmentPermission struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	CanView      bool               `bson:"can_view" json:"can_view"`
	CanDownload  bool               `bson:"can_download" json:"can_download"`
	CanDelete    bool               `bson:"can_delete" json:"can_delete"`
	CanShare     bool               `bson:"can_share" json:"can_share"`
	GrantedBy    string             `bson:"granted_by" json:"granted_by"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// ── Share Links ──

type ShareLink struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	Code         string             `bson:"code" json:"code"`
	CreatedBy    string             `bson:"created_by" json:"created_by"`
	Password     string             `bson:"password,omitempty" json:"password,omitempty"`
	MaxDownloads int                `bson:"max_downloads" json:"max_downloads"`
	DownloadCount int               `bson:"download_count" json:"download_count"`
	ExpiresAt    *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// ── Virus Scan Results ──

type ScanResult struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	Status       string             `bson:"status" json:"status"` // pending, clean, infected, error
	Engine       string             `bson:"engine" json:"engine"`
	Details      string             `bson:"details" json:"details"`
	ScannedAt    time.Time          `bson:"scanned_at" json:"scanned_at"`
}

// ── Previews / Thumbnails ──

type AttachmentPreview struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AttachmentID string             `bson:"attachment_id" json:"attachment_id"`
	PreviewType  string             `bson:"preview_type" json:"preview_type"` // thumbnail, preview, icon
	Width        int                `bson:"width" json:"width"`
	Height       int                `bson:"height" json:"height"`
	URL          string             `bson:"url" json:"url"`
	StoragePath  string             `bson:"storage_path" json:"storage_path"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// ── Quota & Stats ──

type UserQuota struct {
	UserID     string `json:"user_id"`
	UsedBytes  int64  `json:"used_bytes"`
	MaxBytes   int64  `json:"max_bytes"`
	FileCount  int64  `json:"file_count"`
	MaxFiles   int64  `json:"max_files"`
}

type AttachmentStats struct {
	TotalFiles      int64            `json:"total_files"`
	TotalSize       int64            `json:"total_size"`
	ByType          map[string]int64 `json:"by_type"`
	ByStatus        map[string]int64 `json:"by_status"`
	RecentUploads   int64            `json:"recent_uploads_24h"`
}

type WorkspaceStats struct {
	WorkspaceID  string `json:"workspace_id"`
	TotalFiles   int64  `json:"total_files"`
	TotalSize    int64  `json:"total_size"`
	UserCount    int64  `json:"user_count"`
}

// ── Request DTOs ──

type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID string `json:"parent_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type AddTagRequest struct {
	Tag string `json:"tag" binding:"required"`
}

type CreateShareRequest struct {
	SharedWith string `json:"shared_with" binding:"required"`
	Permission string `json:"permission" binding:"required"`
	ExpiresAt  string `json:"expires_at"`
}

type CreateCollectionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type UpdateCollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    *bool  `json:"is_public"`
}

type AddToCollectionRequest struct {
	AttachmentID string `json:"attachment_id" binding:"required"`
}

type SetPermissionRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	CanView     bool   `json:"can_view"`
	CanDownload bool   `json:"can_download"`
	CanDelete   bool   `json:"can_delete"`
	CanShare    bool   `json:"can_share"`
}

type CreateShareLinkRequest struct {
	Password     string `json:"password"`
	MaxDownloads int    `json:"max_downloads"`
	ExpiresAt    string `json:"expires_at"`
}

type BulkDeleteRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

type BulkMoveRequest struct {
	IDs       []string `json:"ids" binding:"required"`
	ChannelID string   `json:"channel_id" binding:"required"`
}

type BulkTagRequest struct {
	IDs []string `json:"ids" binding:"required"`
	Tag string   `json:"tag" binding:"required"`
}

type SearchRequest struct {
	Query       string `form:"q"`
	Type        string `form:"type"`
	ChannelID   string `form:"channel_id"`
	UserID      string `form:"user_id"`
	MinSize     int64  `form:"min_size"`
	MaxSize     int64  `form:"max_size"`
	StartDate   string `form:"start_date"`
	EndDate     string `form:"end_date"`
	Tags        string `form:"tags"`
}

type RenameRequest struct {
	NewName string `json:"new_name" binding:"required"`
}

type MoveRequest struct {
	ChannelID string `json:"channel_id" binding:"required"`
	MessageID string `json:"message_id"`
}

type CloneRequest struct {
	ChannelID string `json:"channel_id" binding:"required"`
	MessageID string `json:"message_id"`
}
