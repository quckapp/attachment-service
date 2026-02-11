package main

import (
	"log"
	"os"

	"attachment-service/internal/api"
	"attachment-service/internal/config"
	"attachment-service/internal/kafka"
	"attachment-service/internal/repository"
	"attachment-service/internal/service"
	"attachment-service/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Initialize MongoDB repository
	repo, err := repository.NewMongoRepository(cfg.MongoDBURL, cfg.DatabaseName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer repo.Close()

	// Initialize extended repository (uses same MongoDB client)
	extRepo := repository.NewExtendedRepository(repo.Client(), cfg.DatabaseName)

	// Initialize storage backend (S3-compatible)
	storageBackend, err := storage.NewS3Storage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize Kafka producer
	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		log.Printf("Warning: Failed to connect Kafka producer: %v", err)
	}
	defer func() {
		if producer != nil {
			producer.Close()
		}
	}()

	// Initialize service
	attachmentService := service.NewAttachmentService(repo, storageBackend, producer, cfg)

	// Setup HTTP server
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	api.RegisterRoutes(router, attachmentService, cfg)
	api.RegisterExtendedRoutes(router, extRepo)
	api.RegisterExtendedRoutes2(router, extRepo, extRepo.Database())

	port := cfg.Port
	log.Printf("Attachment service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
