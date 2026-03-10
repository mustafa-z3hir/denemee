package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              string         `json:"id" gorm:"primaryKey"`
	Email           string         `json:"email" gorm:"uniqueIndex"`
	Username        string         `json:"username" gorm:"uniqueIndex"`
	Password        string         `json:"-"`
	Avatar          string         `json:"avatar"`
	Bio             string         `json:"bio"`
	IsAdmin         bool           `json:"is_admin" gorm:"default:false"`
	IsVerified      bool           `json:"is_verified" gorm:"default:false"`
	EquippedSlots   int            `json:"equipped_slots" gorm:"default:3"`
	MaxSlots        int            `json:"max_slots" gorm:"default:3"`
	AniPoints       int            `json:"ani_points" gorm:"default:0"`
	TotalWatchHours float64        `json:"total_watch_hours" gorm:"default:0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	
	Badges          []Badge          `json:"badges,omitempty"`
}

type Anime struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	TitleEn     string    `json:"title_en"`
	TitleJp     string    `json:"title_jp"`
	Description string    `json:"description"`
	CoverImage  string    `json:"cover_image"`
	BannerImage string    `json:"banner_image"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Season      string    `json:"season"`
	Year        int       `json:"year"`
	Episodes    int       `json:"episodes"`
	Duration    int       `json:"duration"`
	Rating      string    `json:"rating"`
	IMDbID      string    `json:"imdb_id"`
	IMDbRating  float64   `json:"imdb_rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	Genres      []Genre     `json:"genres,omitempty" gorm:"many2many:anime_genres;"`
	Studios     []Studio    `json:"studios,omitempty" gorm:"many2many:anime_studios;"`
	EpisodesList []Episode  `json:"episodes,omitempty" gorm:"foreignKey:AnimeID"`
}

type Episode struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	AnimeID     string    `json:"anime_id"`
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Duration    int       `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	
	Sources     []VideoSource `json:"sources,omitempty"`
	Subtitles   []Subtitle    `json:"subtitles,omitempty"`
}

type VideoSource struct {
	ID        string `json:"id" gorm:"primaryKey"`
	EpisodeID string `json:"episode_id"`
	Label     string `json:"label"`
	URL       string `json:"url"`
	Provider  string `json:"provider"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
}

type Subtitle struct {
	ID        string  `json:"id" gorm:"primaryKey"`
	EpisodeID string  `json:"episode_id"`
	Language  string  `json:"language"`
	URL       string  `json:"url"`
	FansubID  *string `json:"fansub_id,omitempty"`
}

type Genre struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"uniqueIndex"`
}

type Studio struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"uniqueIndex"`
}

type Rating struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	AnimeID   string    `json:"anime_id"`
	UserID    string    `json:"user_id"`
	Score     int       `json:"score"`
	Review    string    `json:"review"`
	IsSpoiler bool      `json:"is_spoiler" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Fansub struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Slug        string    `json:"slug" gorm:"uniqueIndex"`
	Name        string    `json:"name"`
	Logo        string    `json:"logo"`
	Banner      string    `json:"banner"`
	Description string    `json:"description"`
	Website     string    `json:"website"`
	Discord     string    `json:"discord"`
	Twitter     string    `json:"twitter"`
	FoundedAt   *time.Time `json:"founded_at"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	
	Members  []FansubMember  `json:"members,omitempty"`
	Projects []FansubProject `json:"projects,omitempty"`
}

type FansubMember struct {
	ID       string    `json:"id" gorm:"primaryKey"`
	FansubID string    `json:"fansub_id"`
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active" gorm:"default:true"`
	JoinedAt time.Time `json:"joined_at"`
	
	User   User   `json:"user,omitempty"`
	Fansub Fansub `json:"fansub,omitempty"`
}

type FansubProject struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	FansubID  string    `json:"fansub_id"`
	AnimeID   string    `json:"anime_id"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Badge struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Tier        int       `json:"tier"`
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	Gradient    string    `json:"gradient"`
	Animation   string    `json:"animation"`
	Description string    `json:"description"`
	EarnedAt    time.Time `json:"earned_at"`
	EarnedHow   string    `json:"earned_how"`
	IsEquipped  bool      `json:"is_equipped" gorm:"default:false"`
	Slot        int       `json:"slot"`
}

type WatchHistory struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	UserID        string    `json:"user_id"`
	EpisodeID     string    `json:"episode_id"`
	Progress      int       `json:"progress"`
	IsCompleted   bool      `json:"is_completed" gorm:"default:false"`
	WatchCount    int       `json:"watch_count" gorm:"default:1"`
	LastWatchedAt time.Time `json:"last_watched_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type Conversation struct {
	ID            string     `json:"id" gorm:"primaryKey"`
	Type          string     `json:"type"`
	Title         *string    `json:"title,omitempty"`
	Avatar        *string    `json:"avatar,omitempty"`
	LastMessageAt *time.Time `json:"last_message_at"`
	CreatedAt     time.Time  `json:"created_at"`
	
	Participants []ConversationParticipant `json:"participants,omitempty"`
	Messages     []Message                 `json:"messages,omitempty"`
}

type ConversationParticipant struct {
	ID             string     `json:"id" gorm:"primaryKey"`
	ConversationID string     `json:"conversation_id"`
	UserID         string     `json:"user_id"`
	IsAdmin        bool       `json:"is_admin" gorm:"default:false"`
	JoinedAt       time.Time  `json:"joined_at"`
	LastReadAt     *time.Time `json:"last_read_at"`
}

type Message struct {
	ID             string     `json:"id" gorm:"primaryKey"`
	ConversationID string     `json:"conversation_id"`
	SenderID       string     `json:"sender_id"`
	Type           string     `json:"type"`
	Content        string     `json:"content"`
	ReplyToID      *string    `json:"reply_to_id,omitempty"`
	IsEdited       bool       `json:"is_edited" gorm:"default:false"`
	EditedAt       *time.Time `json:"edited_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type Clip struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id"`
	AnimeID     string    `json:"anime_id"`
	EpisodeID   string    `json:"episode_id"`
	StartTime   int       `json:"start_time"`
	EndTime     int       `json:"end_time"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	GIFURL      string    `json:"gif_url"`
	VideoURL    string    `json:"video_url"`
	ViewCount   int       `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
}

type Purchase struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	UserID       string    `json:"user_id"`
	Type         string    `json:"type"`
	ItemID       string    `json:"item_id"`
	ItemName     string    `json:"item_name"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency" gorm:"default:TRY"`
	Points       int       `json:"points"`
	SlotIncrease int       `json:"slot_increase"`
	PurchasedAt  time.Time `json:"purchased_at"`
}

type PointTransaction struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Follow struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	FollowerID  string    `json:"follower_id"`
	FollowingID string    `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Report struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	ReporterID  string     `json:"reporter_id"`
	TargetType  string     `json:"target_type"`
	TargetID    string     `json:"target_id"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Status      string     `json:"status" gorm:"default:pending"`
	ResolvedBy  *string    `json:"resolved_by,omitempty"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}