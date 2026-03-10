package repository

import (
	"context"
	"time"

	"aniverse/internal/domain"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{DB: db, Redis: redis}
}

// User methods
func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).Preload("Badges").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).Preload("Badges", "is_equipped = ?", true).First(&user, "id = ?", id).Error
	return &user, err
}

func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}

// Anime methods
func (r *Repository) GetAnimes(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]domain.Anime, error) {
	var animes []domain.Anime
	query := r.DB.WithContext(ctx).Preload("Genres").Preload("Studios")
	
	if title, ok := filters["title"]; ok {
		query = query.Where("title ILIKE ? OR title_en ILIKE ?", "%"+title.(string)+"%", "%"+title.(string)+"%")
	}
	if genre, ok := filters["genre"]; ok {
		query = query.Joins("JOIN anime_genres ON animes.id = anime_genres.anime_id").
			Joins("JOIN genres ON anime_genres.genre_id = genres.id").
			Where("genres.name = ?", genre)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if year, ok := filters["year"]; ok {
		query = query.Where("year = ?", year)
	}
	
	err := query.Limit(limit).Offset(offset).Find(&animes).Error
	return animes, err
}

func (r *Repository) GetAnimeByID(ctx context.Context, id string) (*domain.Anime, error) {
	var anime domain.Anime
	err := r.DB.WithContext(ctx).Preload("Genres").Preload("Studios").Preload("EpisodesList").First(&anime, "id = ?", id).Error
	return &anime, err
}

func (r *Repository) CreateAnime(ctx context.Context, anime *domain.Anime) error {
	return r.DB.WithContext(ctx).Create(anime).Error
}

// Rating methods
func (r *Repository) CreateOrUpdateRating(ctx context.Context, rating *domain.Rating) error {
	var existing domain.Rating
	err := r.DB.WithContext(ctx).Where("anime_id = ? AND user_id = ?", rating.AnimeID, rating.UserID).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		return r.DB.WithContext(ctx).Create(rating).Error
	}
	
	existing.Score = rating.Score
	existing.Review = rating.Review
	existing.IsSpoiler = rating.IsSpoiler
	return r.DB.WithContext(ctx).Save(&existing).Error
}

func (r *Repository) GetAnimeAverageRating(ctx context.Context, animeID string) (float64, int, error) {
	type Result struct {
		Avg   float64
		Count int64
	}
	var result Result
	err := r.DB.WithContext(ctx).Model(&domain.Rating{}).
		Select("AVG(score) as avg, COUNT(*) as count").
		Where("anime_id = ?", animeID).
		Scan(&result).Error
	return result.Avg, int(result.Count), err
}

// Badge methods
func (r *Repository) CreateBadge(ctx context.Context, badge *domain.Badge) error {
	return r.DB.WithContext(ctx).Create(badge).Error
}

func (r *Repository) GetUserBadges(ctx context.Context, userID string) ([]domain.Badge, error) {
	var badges []domain.Badge
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Order("earned_at DESC").Find(&badges).Error
	return badges, err
}

func (r *Repository) UpdateBadgeEquip(ctx context.Context, userID string, badgeID string, equipped bool, slot int) error {
	return r.DB.WithContext(ctx).Model(&domain.Badge{}).
		Where("id = ? AND user_id = ?", badgeID, userID).
		Updates(map[string]interface{}{"is_equipped": equipped, "slot": slot}).Error
}

// Watch history
func (r *Repository) CreateOrUpdateWatchHistory(ctx context.Context, history *domain.WatchHistory) error {
	var existing domain.WatchHistory
	err := r.DB.WithContext(ctx).Where("user_id = ? AND episode_id = ?", history.UserID, history.EpisodeID).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		return r.DB.WithContext(ctx).Create(history).Error
	}
	
	existing.Progress = history.Progress
	existing.IsCompleted = history.IsCompleted
	existing.WatchCount++
	existing.LastWatchedAt = history.LastWatchedAt
	return r.DB.WithContext(ctx).Save(&existing).Error
}

func (r *Repository) GetUserWatchStats(ctx context.Context, userID string) (float64, int, error) {
	type Result struct {
		TotalHours float64
		Count      int64
	}
	var result Result
	err := r.DB.WithContext(ctx).Model(&domain.WatchHistory{}).
		Select("COALESCE(SUM(progress), 0) / 3600.0 as total_hours, COUNT(*) as count").
		Where("user_id = ?", userID).
		Scan(&result).Error
	return result.TotalHours, int(result.Count), err
}

// Fansub methods
func (r *Repository) GetFansubs(ctx context.Context) ([]domain.Fansub, error) {
	var fansubs []domain.Fansub
	err := r.DB.WithContext(ctx).Where("is_active = ?", true).Find(&fansubs).Error
	return fansubs, err
}

func (r *Repository) GetFansubBySlug(ctx context.Context, slug string) (*domain.Fansub, error) {
	var fansub domain.Fansub
	err := r.DB.WithContext(ctx).Preload("Members.User").Preload("Projects.Anime").First(&fansub, "slug = ?", slug).Error
	return &fansub, err
}

func (r *Repository) CreateFansub(ctx context.Context, fansub *domain.Fansub) error {
	return r.DB.WithContext(ctx).Create(fansub).Error
}

// DM methods
func (r *Repository) GetConversations(ctx context.Context, userID string) ([]domain.Conversation, error) {
	var conversations []domain.Conversation
	err := r.DB.WithContext(ctx).
		Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversation_participants.user_id = ?", userID).
		Preload("Participants.User").
		Order("last_message_at DESC NULLS LAST").
		Find(&conversations).Error
	return conversations, err
}

func (r *Repository) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]domain.Message, error) {
	var messages []domain.Message
	err := r.DB.WithContext(ctx).
		Where("conversation_id = ? AND deleted_at IS NULL", conversationID).
		Preload("Sender").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *Repository) CreateMessage(ctx context.Context, message *domain.Message) error {
	return r.DB.WithContext(ctx).Create(message).Error
}

func (r *Repository) UpdateConversationLastMessage(ctx context.Context, conversationID string, t time.Time) error {
	return r.DB.WithContext(ctx).Model(&domain.Conversation{}).
		Where("id = ?", conversationID).
		Update("last_message_at", t).Error
}

// Shop methods
func (r *Repository) CreatePurchase(ctx context.Context, purchase *domain.Purchase) error {
	return r.DB.WithContext(ctx).Create(purchase).Error
}

func (r *Repository) GetUserPurchases(ctx context.Context, userID string) ([]domain.Purchase, error) {
	var purchases []domain.Purchase
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Order("purchased_at DESC").Find(&purchases).Error
	return purchases, err
}

// Point methods
func (r *Repository) AddPoints(ctx context.Context, userID string, points int, description string) error {
	tx := r.DB.WithContext(ctx).Begin()
	
	if err := tx.Model(&domain.User{}).Where("id = ?", userID).Update("ani_points", gorm.Expr("ani_points + ?", points)).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	transaction := &domain.PointTransaction{
		UserID:      userID,
		Type:        "earned_ad",
		Amount:      points,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit().Error
}

func (r *Repository) GetPointHistory(ctx context.Context, userID string) ([]domain.PointTransaction, error) {
	var transactions []domain.PointTransaction
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}