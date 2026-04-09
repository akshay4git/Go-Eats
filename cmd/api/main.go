package main

import (
	"log"
	"os"

	"github.com/Ayocodes24/GO-Eats/cmd/api/middleware"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/handler"
	annoucements "github.com/Ayocodes24/GO-Eats/pkg/handler/announcements"
	crt "github.com/Ayocodes24/GO-Eats/pkg/handler/cart"
	delv "github.com/Ayocodes24/GO-Eats/pkg/handler/delivery"
	notify "github.com/Ayocodes24/GO-Eats/pkg/handler/notification"
	"github.com/Ayocodes24/GO-Eats/pkg/handler/restaurant"
	revw "github.com/Ayocodes24/GO-Eats/pkg/handler/review"
	"github.com/Ayocodes24/GO-Eats/pkg/handler/user"
	"github.com/Ayocodes24/GO-Eats/pkg/nats"
	"github.com/Ayocodes24/GO-Eats/pkg/service/announcements"
	"github.com/Ayocodes24/GO-Eats/pkg/service/cart_order"
	"github.com/Ayocodes24/GO-Eats/pkg/service/delivery"
	"github.com/Ayocodes24/GO-Eats/pkg/service/notification"
	restro "github.com/Ayocodes24/GO-Eats/pkg/service/restaurant"
	"github.com/Ayocodes24/GO-Eats/pkg/service/review"
	usr "github.com/Ayocodes24/GO-Eats/pkg/service/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/Ayocodes24/GO-Eats/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func runMigrations() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Build from individual env vars (local dev)
		databaseURL = "postgres://" + os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") +
			"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	}

	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		log.Fatalf("Migration init error: %s", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %s", err)
	}
	log.Println("Migrations applied successfully")
}

// @title 			GO-EATS API
// @version			1.0
// @description 	A food delivery backend system with user management, restaurant operations, cart/order processing, delivery management with 2FA and real-time notifications.
// @contact.name 	Akshay Sharma
// @contact.url     github.com/akshay4git/Go-Eats

// @host 			localhost:8080
// @BasePath 		/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token


func main() {
	// Load .env for local dev — on Railway env vars are injected directly
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	env := os.Getenv("APP_ENV")

	// Run migrations before anything else
	runMigrations()

	db := database.New()

	// Initialize Validator
	validate := validator.New()

	// Connect NATS
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://127.0.0.1:4222" //fallback for lacal dev
	}

	natServer, err := nats.NewNATS(natsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}

	// WebSocket Clients
	wsClients := make(map[string]*websocket.Conn)

	s := handler.NewServer(db, true)

	// Middlewares List
	middlewares := []gin.HandlerFunc{middleware.AuthMiddleware()}

	// User
	userService := usr.NewUserService(db, env)
	user.NewUserHandler(s, "/user", userService, validate)

	// Reviews
	reviewService := review.NewReviewService(db, env)
	revw.NewReviewProtectedHandler(s, "/review", reviewService, middlewares, validate)

	// Restaurant
	restaurantService := restro.NewRestaurantService(db, env)
	restaurant.NewRestaurantHandler(s, "/restaurant", restaurantService)

	// Cart
	cartService := cart_order.NewCartService(db, env, natServer)
	crt.NewCartHandler(s, "/cart", cartService, middlewares, validate)

	// Delivery
	deliveryService := delivery.NewDeliveryService(db, env, natServer)
	delv.NewDeliveryHandler(s, "/delivery", deliveryService, middlewares, validate)

	// Events/Announcements
	announceService := announcements.NewAnnouncementService(db, env)
	annoucements.NewAnnouncementHandler(s, "/announcements", announceService, middlewares, validate)

	// Notification
	notifyService := notification.NewNotificationService(db, env, natServer)
	_ = notifyService.SubscribeNewOrders(wsClients)
	_ = notifyService.SubscribeOrderStatus(wsClients)
	notify.NewNotifyHandler(s, "/notify", notifyService, middlewares, validate, wsClients)

	s.Gin.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	s.Gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Fatal(s.Run())
}
