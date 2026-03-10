package handler

import (
	"strconv"

	"aniverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// Auth
func (h *Handler) Register(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	
	resp, err := h.svc.Register(c.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(resp)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	
	resp, err := h.svc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(resp)
}

// Anime
func (h *Handler) GetAnimes(c *fiber.Ctx) error {
	filters := make(map[string]interface{})
	
	if title := c.Query("title"); title != "" {
		filters["title"] = title
	}
	if genre := c.Query("genre"); genre != "" {
		filters["genre"] = genre
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	
	animes, err := h.svc.GetAnimes(c.Context(), filters, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(animes)
}

func (h *Handler) GetAnime(c *fiber.Ctx) error {
	id := c.Params("id")
	anime, err := h.svc.GetAnime(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "anime not found"})
	}
	return c.JSON(anime)
}

func (h *Handler) GetEpisodes(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

// User
func (h *Handler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	user, err := h.svc.GetUser(c.Context(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}

func (h *Handler) GetUserProfile(c *fiber.Ctx) error {
	username := c.Params("username")
	return c.JSON(fiber.Map{"username": username})
}

func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func (h *Handler) FollowUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

// Rating
func (h *Handler) RateAnime(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	animeID := c.Params("id")
	
	type Request struct {
		Score     int    `json:"score"`
		Review    string `json:"review"`
		IsSpoiler bool   `json:"is_spoiler"`
	}
	
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	
	if err := h.svc.RateAnime(c.Context(), userID, animeID, req.Score, req.Review, req.IsSpoiler); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(fiber.Map{"message": "rated successfully"})
}

func (h *Handler) GetRatings(c *fiber.Ctx) error {
	animeID := c.Params("id")
	ratings, err := h.svc.GetAnimeRatings(c.Context(), animeID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ratings)
}

// Watch
func (h *Handler) RecordWatch(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	episodeID := c.Params("id")
	
	type Request struct {
		Progress    int  `json:"progress"`
		IsCompleted bool `json:"is_completed"`
	}
	
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	
	if err := h.svc.RecordWatch(c.Context(), userID, episodeID, req.Progress, req.IsCompleted); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(fiber.Map{"message": "watch recorded"})
}

func (h *Handler) GetWatchHistory(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func (h *Handler) GetContinueWatching(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

// Fansub
func (h *Handler) GetFansubs(c *fiber.Ctx) error {
	fansubs, err := h.svc.GetFansubs(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fansubs)
}

func (h *Handler) GetFansub(c *fiber.Ctx) error {
	slug := c.Params("slug")
	fansub, err := h.svc.GetFansub(c.Context(), slug)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "fansub not found"})
	}
	return c.JSON(fansub)
}

func (h *Handler) CreateFansub(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	
	type Request struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	
	fansub, err := h.svc.CreateFansub(c.Context(), userID, req.Name, req.Slug, req.Description)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.Status(201).JSON(fansub)
}

func (h *Handler) RequestJoinFansub(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

// Badges
func (h *Handler) GetMyBadges(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	badges, err := h.svc.GetUserBadges(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(badges)
}

func (h *Handler) CheckNewBadges(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	h.svc.CheckAndAwardBadges(c.Context(), userID)
	return c.JSON(fiber.Map{"message": "checked"})
}

func (h *Handler) EquipBadge(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

// Shop
func (h *Handler) GetShopItems(c *fiber.Ctx) error {
	items := []map[string]interface{}{
		{"id": "slot_1", "type": "slot", "name": "Extra Slot +1", "price": 30.0, "points": 300},
		{"id": "theme_dark", "type": "theme", "name": "Dark Theme", "price": 15.0, "points": 100},
		{"id": "frame_gold", "type": "frame", "name": "Gold Frame", "price": 9.0, "points": 60},
	}
	return c.JSON(items)
}

func (h *Handler) PurchaseItem(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

// Points
func (h *Handler) GetPointsBalance(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"balance": 0})
}

func (h *Handler) EarnPoints(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

// DM
func (h *Handler) GetConversations(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	conversations, err := h.svc.GetConversations(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(conversations)
}

func (h *Handler) CreateConversation(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func (h *Handler) GetMessages(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func (h *Handler) SendMessage(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

// Clips
func (h *Handler) CreateClip(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func (h *Handler) GetClip(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

// Admin
func (h *Handler) AdminDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"totalUsers":      100,
		"activeToday":     50,
		"totalAnime":      500,
		"pendingReports":  5,
	})
}

func (h *Handler) AdminGetUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func (h *Handler) AdminBanUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}