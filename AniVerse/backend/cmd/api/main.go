package main

import (
	"log"

	"aniverse/internal/config"
	"aniverse/internal/handler"
	"aniverse/internal/middleware"
	"aniverse/internal/repository"
	"aniverse/internal/service"
	"aniverse/internal/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	redis := config.ConnectRedis(cfg)
	repo := repository.NewRepository(db, redis)
	svc := service.NewService(repo, cfg)
	h := handler.NewHandler(svc)
	ws := websocket.NewHub(repo)

	go ws.Run()

	app := fiber.New(fiber.Config{
		AppName: "AniVerse API",
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// --- BURAYI EKLEDİK: ANA SAYFA TEST ROTASI ---
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app":     "AniVerse API",
			"version": "v1.0.0",
			"status":  "Server is running!",
		})
	})
	// --------------------------------------------

	api := app.Group("/api/v1")

	// Public routes
	api.Post("/auth/register", h.Register)
	api.Post("/auth/login", h.Login)

	api.Get("/anime", h.GetAnimes)
	api.Get("/anime/:id", h.GetAnime)
	api.Get("/anime/:id/episodes", h.GetEpisodes)
	api.Get("/anime/:id/ratings", h.GetRatings)

	api.Get("/fansubs", h.GetFansubs)
	api.Get("/fansubs/:slug", h.GetFansub)

	// Protected routes
	protected := api.Group("/", middleware.Auth(cfg.JWTSecret))

	protected.Get("/user/me", h.GetMe)
	protected.Put("/user/me", h.UpdateProfile)
	protected.Get("/user/:username", h.GetUserProfile)
	protected.Post("/user/follow/:id", h.FollowUser)

	protected.Post("/anime/:id/rate", h.RateAnime)
	protected.Post("/anime/:id/watch", h.RecordWatch)

	protected.Get("/history", h.GetWatchHistory)
	protected.Get("/history/continue", h.GetContinueWatching)

	protected.Post("/fansubs", h.CreateFansub)
	protected.Post("/fansubs/:id/join", h.RequestJoinFansub)

	protected.Get("/badges", h.GetMyBadges)
	protected.Get("/badges/check", h.CheckNewBadges)
	protected.Post("/badges/equip", h.EquipBadge)

	protected.Get("/shop/items", h.GetShopItems)
	protected.Post("/shop/purchase", h.PurchaseItem)

	protected.Get("/points/balance", h.GetPointsBalance)
	protected.Post("/points/earn", h.EarnPoints)

	protected.Get("/conversations", h.GetConversations)
	protected.Post("/conversations", h.CreateConversation)
	protected.Get("/conversations/:id/messages", h.GetMessages)
	protected.Post("/conversations/:id/messages", h.SendMessage)

	protected.Post("/clips", h.CreateClip)
	protected.Get("/clips/:id", h.GetClip)

	// Admin routes
	admin := protected.Group("/admin", middleware.Auth(cfg.JWTSecret), middleware.AdminOnly())
	admin.Get("/dashboard", h.AdminDashboard)
	admin.Get("/users", h.AdminGetUsers)
	admin.Put("/users/:id/ban", h.AdminBanUser)

	app.Get("/ws", middleware.Auth(cfg.JWTSecret), ws.HandleWebSocket)

	log.Fatal(app.Listen(":8080"))
}
