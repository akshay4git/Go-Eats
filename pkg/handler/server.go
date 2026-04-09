package handler

import (
	"log/slog"
	"os"
	"time"

	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type Server struct {
	Gin     *gin.Engine
	db      database.Database
	Storage storage.ImageStorage
}

func NewServer(db database.Database, setLog bool) *Server {
	ginEngine := gin.New()

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}

	// Rate limiter: 100 requests per minute per IP
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	// Setting Logger, CORS, Rate Limiter & MultipartMemory
	if setLog {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		ginEngine.Use(sloggin.New(logger))
	}
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(cors.New(corsConfig))
	ginEngine.Use(ginlimiter.NewMiddleware(instance))
	ginEngine.MaxMultipartMemory = 8 << 20 // 8 MB

	localStoragePath := os.Getenv("LOCAL_STORAGE_PATH")
	if len(localStoragePath) > 0 {
		ginEngine.Static("/uploads", localStoragePath)
	}

	return &Server{
		Gin:     ginEngine,
		db:      db,
		Storage: storage.CreateImageStorage(os.Getenv("STORAGE_TYPE")),
	}
}

func (server *Server) Run() error {
	return server.Gin.Run(":8080")
}
