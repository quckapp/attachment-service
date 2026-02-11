package repository

import (
	"context"
	"time"

	"attachment-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ExtendedRepository provides additional MongoDB operations
type ExtendedRepository struct {
	client     *mongo.Client
	db         *mongo.Database
	versions   *mongo.Collection
	comments   *mongo.Collection
	tags       *mongo.Collection
	favorites  *mongo.Collection
	shares     *mongo.Collection
	collections *mongo.Collection
	collItems  *mongo.Collection
	activities *mongo.Collection
	permissions *mongo.Collection
	shareLinks *mongo.Collection
	scans      *mongo.Collection
	previews   *mongo.Collection
	attachments *mongo.Collection
}

func NewExtendedRepository(client *mongo.Client, dbName string) *ExtendedRepository {
	db := client.Database(dbName)
	r := &ExtendedRepository{
		client:      client,
		db:          db,
		versions:    db.Collection("attachment_versions"),
		comments:    db.Collection("attachment_comments"),
		tags:        db.Collection("attachment_tags"),
		favorites:   db.Collection("attachment_favorites"),
		shares:      db.Collection("attachment_shares"),
		collections: db.Collection("attachment_collections"),
		collItems:   db.Collection("collection_items"),
		activities:  db.Collection("attachment_activities"),
		permissions: db.Collection("attachment_permissions"),
		shareLinks:  db.Collection("share_links"),
		scans:       db.Collection("scan_results"),
		previews:    db.Collection("attachment_previews"),
		attachments: db.Collection("attachments"),
	}

	ctx := context.Background()
	// Create indexes
	r.versions.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "attachment_id", Value: 1}}},
		{Keys: bson.D{{Key: "attachment_id", Value: 1}, {Key: "version_num", Value: -1}}},
	})
	r.comments.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "attachment_id", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
	})
	r.tags.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "attachment_id", Value: 1}}},
		{Keys: bson.D{{Key: "tag", Value: 1}}},
		{Keys: bson.D{{Key: "attachment_id", Value: 1}, {Key: "tag", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	r.favorites.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "attachment_id", Value: 1}, {Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	r.shares.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "attachment_id", Value: 1}}},
		{Keys: bson.D{{Key: "shared_with", Value: 1}}},
	})
	r.collections.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "workspace_id", Value: 1}}},
		{Keys: bson.D{{Key: "created_by", Value: 1}}},
	})
	r.collItems.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "collection_id", Value: 1}}},
		{Keys: bson.D{{Key: "collection_id", Value: 1}, {Key: "attachment_id", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	r.activities.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "attachment_id", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	r.permissions.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "attachment_id", Value: 1}}},
		{Keys: bson.D{{Key: "attachment_id", Value: 1}, {Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	r.shareLinks.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "attachment_id", Value: 1}}},
		{Keys: bson.D{{Key: "code", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	r.scans.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "attachment_id", Value: 1}},
	})
	r.previews.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "attachment_id", Value: 1}},
	})

	return r
}

// ── Version Operations ──

func (r *ExtendedRepository) CreateVersion(ctx context.Context, v *models.AttachmentVersion) error {
	v.CreatedAt = time.Now()
	result, err := r.versions.InsertOne(ctx, v)
	if err != nil {
		return err
	}
	v.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ExtendedRepository) ListVersions(ctx context.Context, attachmentID string) ([]*models.AttachmentVersion, error) {
	cursor, err := r.versions.Find(ctx, bson.M{"attachment_id": attachmentID},
		options.Find().SetSort(bson.D{{Key: "version_num", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var versions []*models.AttachmentVersion
	if err := cursor.All(ctx, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

func (r *ExtendedRepository) GetVersionByNum(ctx context.Context, attachmentID string, versionNum int) (*models.AttachmentVersion, error) {
	var v models.AttachmentVersion
	err := r.versions.FindOne(ctx, bson.M{"attachment_id": attachmentID, "version_num": versionNum}).Decode(&v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *ExtendedRepository) GetLatestVersionNum(ctx context.Context, attachmentID string) (int, error) {
	var v models.AttachmentVersion
	err := r.versions.FindOne(ctx, bson.M{"attachment_id": attachmentID},
		options.FindOne().SetSort(bson.D{{Key: "version_num", Value: -1}})).Decode(&v)
	if err == mongo.ErrNoDocuments {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return v.VersionNum, nil
}

func (r *ExtendedRepository) DeleteVersion(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.versions.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// ── Comment Operations ──

func (r *ExtendedRepository) CreateComment(ctx context.Context, c *models.AttachmentComment) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	result, err := r.comments.InsertOne(ctx, c)
	if err != nil {
		return err
	}
	c.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ExtendedRepository) GetComment(ctx context.Context, id string) (*models.AttachmentComment, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var c models.AttachmentComment
	err = r.comments.FindOne(ctx, bson.M{"_id": objID}).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ExtendedRepository) ListComments(ctx context.Context, attachmentID string, limit, offset int) ([]*models.AttachmentComment, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.comments.Find(ctx, bson.M{"attachment_id": attachmentID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var comments []*models.AttachmentComment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *ExtendedRepository) UpdateComment(ctx context.Context, id, content string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.comments.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{
		"content":    content,
		"is_edited":  true,
		"updated_at": time.Now(),
	}})
	return err
}

func (r *ExtendedRepository) DeleteComment(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.comments.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// ── Tag Operations ──

func (r *ExtendedRepository) AddTag(ctx context.Context, t *models.AttachmentTag) error {
	t.CreatedAt = time.Now()
	result, err := r.tags.InsertOne(ctx, t)
	if err != nil {
		return err
	}
	t.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ExtendedRepository) RemoveTag(ctx context.Context, attachmentID, tag string) error {
	_, err := r.tags.DeleteOne(ctx, bson.M{"attachment_id": attachmentID, "tag": tag})
	return err
}

func (r *ExtendedRepository) ListTags(ctx context.Context, attachmentID string) ([]*models.AttachmentTag, error) {
	cursor, err := r.tags.Find(ctx, bson.M{"attachment_id": attachmentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tags []*models.AttachmentTag
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *ExtendedRepository) SearchByTag(ctx context.Context, tag string, limit, offset int) ([]*models.AttachmentTag, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.tags.Find(ctx, bson.M{"tag": tag}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tags []*models.AttachmentTag
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// ── Favorite Operations ──

func (r *ExtendedRepository) AddFavorite(ctx context.Context, f *models.AttachmentFavorite) error {
	f.CreatedAt = time.Now()
	result, err := r.favorites.InsertOne(ctx, f)
	if err != nil {
		return err
	}
	f.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ExtendedRepository) RemoveFavorite(ctx context.Context, attachmentID, userID string) error {
	_, err := r.favorites.DeleteOne(ctx, bson.M{"attachment_id": attachmentID, "user_id": userID})
	return err
}

func (r *ExtendedRepository) ListFavorites(ctx context.Context, userID string, limit, offset int) ([]*models.AttachmentFavorite, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.favorites.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var favs []*models.AttachmentFavorite
	if err := cursor.All(ctx, &favs); err != nil {
		return nil, err
	}
	return favs, nil
}

func (r *ExtendedRepository) IsFavorited(ctx context.Context, attachmentID, userID string) (bool, error) {
	count, err := r.favorites.CountDocuments(ctx, bson.M{"attachment_id": attachmentID, "user_id": userID})
	return count > 0, err
}

// ── Share Operations ──

func (r *ExtendedRepository) CreateShare(ctx context.Context, s *models.AttachmentShare) error {
	s.CreatedAt = time.Now()
	result, err := r.shares.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	s.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ExtendedRepository) ListShares(ctx context.Context, attachmentID string) ([]*models.AttachmentShare, error) {
	cursor, err := r.shares.Find(ctx, bson.M{"attachment_id": attachmentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var shares []*models.AttachmentShare
	if err := cursor.All(ctx, &shares); err != nil {
		return nil, err
	}
	return shares, nil
}

func (r *ExtendedRepository) ListSharedWith(ctx context.Context, userID string, limit, offset int) ([]*models.AttachmentShare, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.shares.Find(ctx, bson.M{"shared_with": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var shares []*models.AttachmentShare
	if err := cursor.All(ctx, &shares); err != nil {
		return nil, err
	}
	return shares, nil
}

func (r *ExtendedRepository) DeleteShare(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.shares.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// ── Collection Operations ──

func (r *ExtendedRepository) CreateCollection(ctx context.Context, c *models.AttachmentCollection) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	result, err := r.collections.InsertOne(ctx, c)
	if err != nil {
		return err
	}
	c.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ExtendedRepository) GetCollection(ctx context.Context, id string) (*models.AttachmentCollection, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var c models.AttachmentCollection
	err = r.collections.FindOne(ctx, bson.M{"_id": objID}).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ExtendedRepository) ListCollections(ctx context.Context, workspaceID string, limit, offset int) ([]*models.AttachmentCollection, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.collections.Find(ctx, bson.M{"workspace_id": workspaceID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var colls []*models.AttachmentCollection
	if err := cursor.All(ctx, &colls); err != nil {
		return nil, err
	}
	return colls, nil
}

func (r *ExtendedRepository) UpdateCollection(ctx context.Context, id string, update bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update["updated_at"] = time.Now()
	_, err = r.collections.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	return err
}

func (r *ExtendedRepository) DeleteCollection(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collections.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	// Delete all items in this collection
	_, err = r.collItems.DeleteMany(ctx, bson.M{"collection_id": id})
	return err
}

func (r *ExtendedRepository) AddToCollection(ctx context.Context, item *models.CollectionItem) error {
	item.AddedAt = time.Now()
	result, err := r.collItems.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	item.ID = result.InsertedID.(primitive.ObjectID)
	// Increment item count
	objID, _ := primitive.ObjectIDFromHex(item.CollectionID)
	_, _ = r.collections.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$inc": bson.M{"item_count": 1}})
	return nil
}

func (r *ExtendedRepository) RemoveFromCollection(ctx context.Context, collectionID, attachmentID string) error {
	_, err := r.collItems.DeleteOne(ctx, bson.M{"collection_id": collectionID, "attachment_id": attachmentID})
	if err != nil {
		return err
	}
	objID, _ := primitive.ObjectIDFromHex(collectionID)
	_, _ = r.collections.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$inc": bson.M{"item_count": -1}})
	return nil
}

func (r *ExtendedRepository) ListCollectionItems(ctx context.Context, collectionID string, limit, offset int) ([]*models.CollectionItem, error) {
	opts := options.Find().SetSort(bson.D{{Key: "position", Value: 1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.collItems.Find(ctx, bson.M{"collection_id": collectionID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var items []*models.CollectionItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// ── Activity Operations ──

func (r *ExtendedRepository) LogActivity(ctx context.Context, a *models.AttachmentActivity) error {
	a.CreatedAt = time.Now()
	_, err := r.activities.InsertOne(ctx, a)
	return err
}

func (r *ExtendedRepository) ListActivity(ctx context.Context, attachmentID string, limit, offset int) ([]*models.AttachmentActivity, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.activities.Find(ctx, bson.M{"attachment_id": attachmentID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var acts []*models.AttachmentActivity
	if err := cursor.All(ctx, &acts); err != nil {
		return nil, err
	}
	return acts, nil
}

func (r *ExtendedRepository) ListUserActivity(ctx context.Context, userID string, limit, offset int) ([]*models.AttachmentActivity, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.activities.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var acts []*models.AttachmentActivity
	if err := cursor.All(ctx, &acts); err != nil {
		return nil, err
	}
	return acts, nil
}

// ── Permission Operations ──

func (r *ExtendedRepository) SetPermission(ctx context.Context, p *models.AttachmentPermission) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	filter := bson.M{"attachment_id": p.AttachmentID, "user_id": p.UserID}
	update := bson.M{"$set": bson.M{
		"can_view":     p.CanView,
		"can_download": p.CanDownload,
		"can_delete":   p.CanDelete,
		"can_share":    p.CanShare,
		"granted_by":   p.GrantedBy,
		"updated_at":   time.Now(),
	}, "$setOnInsert": bson.M{
		"created_at": time.Now(),
	}}
	opts := options.Update().SetUpsert(true)
	_, err := r.permissions.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *ExtendedRepository) GetPermission(ctx context.Context, attachmentID, userID string) (*models.AttachmentPermission, error) {
	var p models.AttachmentPermission
	err := r.permissions.FindOne(ctx, bson.M{"attachment_id": attachmentID, "user_id": userID}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ExtendedRepository) ListPermissions(ctx context.Context, attachmentID string) ([]*models.AttachmentPermission, error) {
	cursor, err := r.permissions.Find(ctx, bson.M{"attachment_id": attachmentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var perms []*models.AttachmentPermission
	if err := cursor.All(ctx, &perms); err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *ExtendedRepository) DeletePermission(ctx context.Context, attachmentID, userID string) error {
	_, err := r.permissions.DeleteOne(ctx, bson.M{"attachment_id": attachmentID, "user_id": userID})
	return err
}

// ── Share Link Operations ──

func (r *ExtendedRepository) CreateShareLink(ctx context.Context, link *models.ShareLink) error {
	link.CreatedAt = time.Now()
	link.IsActive = true
	result, err := r.shareLinks.InsertOne(ctx, link)
	if err != nil {
		return err
	}
	link.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ExtendedRepository) GetShareLinkByCode(ctx context.Context, code string) (*models.ShareLink, error) {
	var link models.ShareLink
	err := r.shareLinks.FindOne(ctx, bson.M{"code": code, "is_active": true}).Decode(&link)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *ExtendedRepository) ListShareLinks(ctx context.Context, attachmentID string) ([]*models.ShareLink, error) {
	cursor, err := r.shareLinks.Find(ctx, bson.M{"attachment_id": attachmentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var links []*models.ShareLink
	if err := cursor.All(ctx, &links); err != nil {
		return nil, err
	}
	return links, nil
}

func (r *ExtendedRepository) IncrementShareLinkDownloads(ctx context.Context, code string) error {
	_, err := r.shareLinks.UpdateOne(ctx, bson.M{"code": code}, bson.M{"$inc": bson.M{"download_count": 1}})
	return err
}

func (r *ExtendedRepository) DeactivateShareLink(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.shareLinks.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"is_active": false}})
	return err
}

// ── Scan Operations ──

func (r *ExtendedRepository) CreateScanResult(ctx context.Context, s *models.ScanResult) error {
	s.ScannedAt = time.Now()
	result, err := r.scans.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	s.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ExtendedRepository) GetScanResult(ctx context.Context, attachmentID string) (*models.ScanResult, error) {
	var s models.ScanResult
	err := r.scans.FindOne(ctx, bson.M{"attachment_id": attachmentID},
		options.FindOne().SetSort(bson.D{{Key: "scanned_at", Value: -1}})).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// ── Preview Operations ──

func (r *ExtendedRepository) CreatePreview(ctx context.Context, p *models.AttachmentPreview) error {
	p.CreatedAt = time.Now()
	result, err := r.previews.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	p.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ExtendedRepository) ListPreviews(ctx context.Context, attachmentID string) ([]*models.AttachmentPreview, error) {
	cursor, err := r.previews.Find(ctx, bson.M{"attachment_id": attachmentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var previews []*models.AttachmentPreview
	if err := cursor.All(ctx, &previews); err != nil {
		return nil, err
	}
	return previews, nil
}

// ── Stats & Search Operations ──

func (r *ExtendedRepository) GetAttachmentStats(ctx context.Context, workspaceID string) (*models.AttachmentStats, error) {
	stats := &models.AttachmentStats{
		ByType:   make(map[string]int64),
		ByStatus: make(map[string]int64),
	}

	filter := bson.M{"workspace_id": workspaceID, "status": bson.M{"$ne": "deleted"}}

	totalCount, err := r.attachments.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	stats.TotalFiles = totalCount

	// Aggregate by type
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$type"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
		}}},
	}
	cursor, err := r.attachments.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			ID        string `bson:"_id"`
			Count     int64  `bson:"count"`
			TotalSize int64  `bson:"total_size"`
		}
		if err := cursor.Decode(&result); err == nil {
			stats.ByType[result.ID] = result.Count
			stats.TotalSize += result.TotalSize
		}
	}

	// Recent uploads (last 24h)
	dayAgo := time.Now().Add(-24 * time.Hour)
	recentCount, _ := r.attachments.CountDocuments(ctx, bson.M{
		"workspace_id": workspaceID,
		"created_at":   bson.M{"$gte": dayAgo},
	})
	stats.RecentUploads = recentCount

	return stats, nil
}

func (r *ExtendedRepository) GetUserQuota(ctx context.Context, userID string) (*models.UserQuota, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userID, "status": bson.M{"$ne": "deleted"}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "used_bytes", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "file_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := r.attachments.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	quota := &models.UserQuota{
		UserID:   userID,
		MaxBytes: 5 * 1024 * 1024 * 1024, // 5GB default
		MaxFiles: 10000,
	}
	if cursor.Next(ctx) {
		var result struct {
			UsedBytes int64 `bson:"used_bytes"`
			FileCount int64 `bson:"file_count"`
		}
		if err := cursor.Decode(&result); err == nil {
			quota.UsedBytes = result.UsedBytes
			quota.FileCount = result.FileCount
		}
	}

	return quota, nil
}

func (r *ExtendedRepository) SearchAttachments(ctx context.Context, workspaceID string, query string, fileType string, limit, offset int) ([]*models.Attachment, error) {
	filter := bson.M{
		"workspace_id": workspaceID,
		"status":       bson.M{"$ne": "deleted"},
	}

	if query != "" {
		filter["original_name"] = bson.M{"$regex": query, "$options": "i"}
	}
	if fileType != "" {
		filter["type"] = fileType
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.attachments.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var attachments []*models.Attachment
	if err := cursor.All(ctx, &attachments); err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *ExtendedRepository) GetRecentAttachments(ctx context.Context, userID string, limit int) ([]*models.Attachment, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	cursor, err := r.attachments.Find(ctx, bson.M{
		"user_id": userID,
		"status":  bson.M{"$ne": "deleted"},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var attachments []*models.Attachment
	if err := cursor.All(ctx, &attachments); err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *ExtendedRepository) GetWorkspaceStats(ctx context.Context, workspaceID string) (*models.WorkspaceStats, error) {
	filter := bson.M{"workspace_id": workspaceID, "status": bson.M{"$ne": "deleted"}}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total_files", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "users", Value: bson.D{{Key: "$addToSet", Value: "$user_id"}}},
		}}},
	}

	cursor, err := r.attachments.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	stats := &models.WorkspaceStats{WorkspaceID: workspaceID}
	if cursor.Next(ctx) {
		var result struct {
			TotalFiles int64    `bson:"total_files"`
			TotalSize  int64    `bson:"total_size"`
			Users      []string `bson:"users"`
		}
		if err := cursor.Decode(&result); err == nil {
			stats.TotalFiles = result.TotalFiles
			stats.TotalSize = result.TotalSize
			stats.UserCount = int64(len(result.Users))
		}
	}

	return stats, nil
}

// ── Bulk Operations ──

func (r *ExtendedRepository) BulkDelete(ctx context.Context, ids []string) error {
	var objIDs []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objIDs = append(objIDs, objID)
	}
	_, err := r.attachments.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": objIDs}}, bson.M{"$set": bson.M{
		"status":     "deleted",
		"updated_at": time.Now(),
	}})
	return err
}

func (r *ExtendedRepository) BulkMove(ctx context.Context, ids []string, channelID string) error {
	var objIDs []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objIDs = append(objIDs, objID)
	}
	_, err := r.attachments.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": objIDs}}, bson.M{"$set": bson.M{
		"channel_id": channelID,
		"updated_at": time.Now(),
	}})
	return err
}

func (r *ExtendedRepository) GetByWorkspaceID(ctx context.Context, workspaceID string, limit, offset int) ([]*models.Attachment, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.attachments.Find(ctx, bson.M{
		"workspace_id": workspaceID,
		"status":       bson.M{"$ne": "deleted"},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var attachments []*models.Attachment
	if err := cursor.All(ctx, &attachments); err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *ExtendedRepository) RenameAttachment(ctx context.Context, id, newName string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.attachments.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{
		"original_name": newName,
		"updated_at":    time.Now(),
	}})
	return err
}

func (r *ExtendedRepository) MoveAttachment(ctx context.Context, id, channelID, messageID string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{
		"channel_id": channelID,
		"updated_at": time.Now(),
	}
	if messageID != "" {
		update["message_id"] = messageID
	}
	_, err = r.attachments.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	return err
}

func (r *ExtendedRepository) CloneAttachment(ctx context.Context, attachment *models.Attachment) error {
	attachment.ID = primitive.NewObjectID()
	attachment.CreatedAt = time.Now()
	attachment.UpdatedAt = time.Now()
	_, err := r.attachments.InsertOne(ctx, attachment)
	return err
}

func (r *ExtendedRepository) AggregateAttachments(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	return r.attachments.Aggregate(ctx, pipeline)
}

func (r *ExtendedRepository) FindAttachmentByID(ctx context.Context, id string, result *models.Attachment) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.attachments.FindOne(ctx, bson.M{"_id": objID}).Decode(result)
}

func (r *ExtendedRepository) Database() *mongo.Database {
	return r.db
}
