package service

import (
	"context"
	"fmt"
	"time"

	"aniverse/internal/config"
	"aniverse/internal/domain"
	"aniverse/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   *repository.Repository
	config *config.Config
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, config: cfg}
}

// Auth
type AuthResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func (s *Service) Register(ctx context.Context, email, username, password string) (*AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	user := &domain.User{
		ID:        generateID(),
		Email:     email,
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	
	s.checkAndAwardBadges(ctx, user.ID)
	
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}
	
	user.Password = ""
	return &AuthResponse{Token: token, User: *user}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}
	
	user.Password = ""
	return &AuthResponse{Token: token, User: *user}, nil
}

func (s *Service) generateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(168 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.config.JWTSecret))
}

// Anime
func (s *Service) GetAnimes(ctx context.Context, filters map[string]interface{}, page, limit int) ([]domain.Anime, error) {
	offset := (page - 1) * limit
	return s.repo.GetAnimes(ctx, filters, limit, offset)
}

func (s *Service) GetAnime(ctx context.Context, id string) (*domain.Anime, error) {
	return s.repo.GetAnimeByID(ctx, id)
}

// Rating
func (s *Service) RateAnime(ctx context.Context, userID, animeID string, score int, review string, isSpoiler bool) error {
	if score < 1 || score > 10 {
		return fmt.Errorf("score must be between 1 and 10")
	}
	
	rating := &domain.Rating{
		ID:        generateID(),
		AnimeID:   animeID,
		UserID:    userID,
		Score:     score,
		Review:    review,
		IsSpoiler: isSpoiler,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return s.repo.CreateOrUpdateRating(ctx, rating)
}

func (s *Service) GetAnimeRatings(ctx context.Context, animeID string) (map[string]interface{}, error) {
	avg, count, err := s.repo.GetAnimeAverageRating(ctx, animeID)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"average": avg,
		"count":   count,
	}, nil
}

// Watch & Badges
func (s *Service) RecordWatch(ctx context.Context, userID, episodeID string, progress int, isCompleted bool) error {
	history := &domain.WatchHistory{
		ID:            generateID(),
		UserID:        userID,
		EpisodeID:     episodeID,
		Progress:      progress,
		IsCompleted:   isCompleted,
		LastWatchedAt: time.Now(),
		CreatedAt:     time.Now(),
	}
	
	if err := s.repo.CreateOrUpdateWatchHistory(ctx, history); err != nil {
		return err
	}
	
	s.checkAndAwardBadges(ctx, userID)
	return nil
}

func (s *Service) checkAndAwardBadges(ctx context.Context, userID string) {
	totalHours, _, _ := s.repo.GetUserWatchStats(ctx, userID)
	
	levels := []struct {
		hours float64
		name  string
		tier  int
	}{
		{5, "Novice", 1},
		{25, "Beginner", 2},
		{100, "Watcher", 3},
		{250, "Viewer", 4},
		{500, "Enthusiast", 5},
		{1000, "Fan", 6},
		{2500, "Otaku", 7},
		{5000, "Master", 8},
		{10000, "Legend", 9},
		{10000, "Immortal", 10},
	}
	
	for _, level := range levels {
		if totalHours >= level.hours {
			existingBadges, _ := s.repo.GetUserBadges(ctx, userID)
			hasBadge := false
			for _, b := range existingBadges {
				if b.Category == level.name {
					hasBadge = true
					break
				}
			}
			
			if !hasBadge {
				badge := &domain.Badge{
					ID:          generateID(),
					UserID:      userID,
					Type:        "watch_level",
					Category:    level.name,
					Tier:        level.tier,
					Name:        level.name,
					Icon:        getLevelIcon(level.name),
					Color:       getLevelColor(level.tier),
					Animation:   getLevelAnimation(level.tier),
					Description: fmt.Sprintf("Watched %.0f hours of anime", level.hours),
					EarnedAt:    time.Now(),
					EarnedHow:   fmt.Sprintf("watched_%.0fh", level.hours),
				}
				s.repo.CreateBadge(ctx, badge)
			}
		}
	}
}

// Fansub
func (s *Service) CreateFansub(ctx context.Context, userID string, name, slug, description string) (*domain.Fansub, error) {
	fansub := &domain.Fansub{
		ID:          generateID(),
		Slug:        slug,
		Name:        name,
		Description: description,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
	
	if err := s.repo.CreateFansub(ctx, fansub); err != nil {
		return nil, err
	}
	
	return fansub, nil
}

func (s *Service) GetFansubs(ctx context.Context) ([]domain.Fansub, error) {
	return s.repo.GetFansubs(ctx)
}

func (s *Service) GetFansub(ctx context.Context, slug string) (*domain.Fansub, error) {
	return s.repo.GetFansubBySlug(ctx, slug)
}

// DM
func (s *Service) GetConversations(ctx context.Context, userID string) ([]domain.Conversation, error) {
	return s.repo.GetConversations(ctx, userID)
}

func (s *Service) SendMessage(ctx context.Context, conversationID, senderID string, content string, msgType string) (*domain.Message, error) {
	message := &domain.Message{
		ID:             generateID(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Type:           msgType,
		Content:        content,
		CreatedAt:      time.Now(),
	}
	
	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return nil, err
	}
	
	s.repo.UpdateConversationLastMessage(ctx, conversationID, time.Now())
	
	return message, nil
}

// User
func (s *Service) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) GetUserBadges(ctx context.Context, userID string) ([]domain.Badge, error) {
	return s.repo.GetUserBadges(ctx, userID)
}

func (s *Service) CheckAndAwardBadges(ctx context.Context, userID string) {
	s.checkAndAwardBadges(ctx, userID)
}

// Helpers
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getLevelIcon(level string) string {
	icons := map[string]string{
		"Novice":     "🌱",
		"Beginner":   "🌿",
		"Watcher":    "🍃",
		"Viewer":     "🌳",
		"Enthusiast": "🌲",
		"Fan":        "🔥",
		"Otaku":      "⭐",
		"Master":     "🌟",
		"Legend":     "👑",
		"Immortal":   "🌌",
	}
	return icons[level]
}

func getLevelColor(tier int) string {
	colors := []string{
		"#6b7280", "#22c55e", "#3b82f6", "#8b5cf6", "#f59e0b",
		"#ef4444", "#ec4899", "#14b8a6", "#fbbf24", "#a855f7",
	}
	if tier > 0 && tier <= len(colors) {
		return colors[tier-1]
	}
	return "#6b7280"
}

func getLevelAnimation(tier int) string {
	if tier <= 3 {
		return "none"
	}
	if tier <= 5 {
		return "pulse"
	}
	if tier <= 7 {
		return "glow"
	}
	if tier <= 9 {
		return "shimmer"
	}
	return "rainbow"
}