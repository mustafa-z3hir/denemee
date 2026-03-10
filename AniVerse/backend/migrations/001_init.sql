CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(500),
    bio TEXT,
    is_admin BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    equipped_slots INTEGER DEFAULT 3,
    max_slots INTEGER DEFAULT 3,
    ani_points INTEGER DEFAULT 0,
    total_watch_hours FLOAT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE animes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    title_en VARCHAR(255),
    title_jp VARCHAR(255),
    description TEXT,
    cover_image VARCHAR(500),
    banner_image VARCHAR(500),
    type VARCHAR(50),
    status VARCHAR(50),
    season VARCHAR(50),
    year INTEGER,
    episodes INTEGER,
    duration INTEGER,
    rating VARCHAR(10),
    imdb_id VARCHAR(20),
    imdb_rating FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE episodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    anime_id UUID REFERENCES animes(id) ON DELETE CASCADE,
    number INTEGER NOT NULL,
    title VARCHAR(255),
    description TEXT,
    thumbnail VARCHAR(500),
    duration INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE video_sources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    episode_id UUID REFERENCES episodes(id) ON DELETE CASCADE,
    label VARCHAR(50),
    url TEXT NOT NULL,
    provider VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE subtitles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    episode_id UUID REFERENCES episodes(id) ON DELETE CASCADE,
    language VARCHAR(10),
    url TEXT NOT NULL,
    fansub_id UUID
);

CREATE TABLE genres (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE anime_genres (
    anime_id UUID REFERENCES animes(id) ON DELETE CASCADE,
    genre_id UUID REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (anime_id, genre_id)
);

CREATE TABLE studios (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE anime_studios (
    anime_id UUID REFERENCES animes(id) ON DELETE CASCADE,
    studio_id UUID REFERENCES studios(id) ON DELETE CASCADE,
    PRIMARY KEY (anime_id, studio_id)
);

CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    anime_id UUID REFERENCES animes(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    score INTEGER CHECK (score >= 1 AND score <= 10),
    review TEXT,
    is_spoiler BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(anime_id, user_id)
);

CREATE TABLE fansubs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    logo VARCHAR(500),
    banner VARCHAR(500),
    description TEXT,
    website VARCHAR(255),
    discord VARCHAR(255),
    twitter VARCHAR(255),
    founded_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE fansub_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    fansub_id UUID REFERENCES fansubs(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(fansub_id, user_id)
);

CREATE TABLE fansub_projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    fansub_id UUID REFERENCES fansubs(id) ON DELETE CASCADE,
    anime_id UUID REFERENCES animes(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'ongoing',
    progress INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE badges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    tier INTEGER DEFAULT 1,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(20),
    gradient VARCHAR(100),
    animation VARCHAR(50),
    description TEXT,
    earned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    earned_how VARCHAR(255),
    is_equipped BOOLEAN DEFAULT FALSE,
    slot INTEGER,
    UNIQUE(user_id, category)
);

CREATE TABLE watch_histories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    episode_id UUID REFERENCES episodes(id) ON DELETE CASCADE,
    progress INTEGER DEFAULT 0,
    is_completed BOOLEAN DEFAULT FALSE,
    watch_count INTEGER DEFAULT 1,
    last_watched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, episode_id)
);

CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(20) DEFAULT 'direct',
    title VARCHAR(255),
    avatar VARCHAR(500),
    last_message_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE conversation_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    is_admin BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_read_at TIMESTAMP,
    UNIQUE(conversation_id, user_id)
);

CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) DEFAULT 'text',
    content TEXT,
    reply_to_id UUID REFERENCES messages(id),
    is_edited BOOLEAN DEFAULT FALSE,
    edited_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE clips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    anime_id UUID REFERENCES animes(id) ON DELETE CASCADE,
    episode_id UUID REFERENCES episodes(id) ON DELETE CASCADE,
    start_time INTEGER NOT NULL,
    end_time INTEGER NOT NULL,
    title VARCHAR(255),
    description TEXT,
    thumbnail VARCHAR(500),
    gif_url VARCHAR(500),
    video_url VARCHAR(500),
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE purchases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    item_id VARCHAR(100) NOT NULL,
    item_name VARCHAR(255),
    price FLOAT,
    currency VARCHAR(10) DEFAULT 'TRY',
    points INTEGER DEFAULT 0,
    slot_increase INTEGER DEFAULT 0,
    purchased_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE point_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    amount INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE follows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    follower_id UUID REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(follower_id, following_id)
);

CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_id UUID REFERENCES users(id) ON DELETE CASCADE,
    target_type VARCHAR(50) NOT NULL,
    target_id UUID NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_animes_status ON animes(status);
CREATE INDEX idx_animes_year ON animes(year);
CREATE INDEX idx_episodes_anime_id ON episodes(anime_id);
CREATE INDEX idx_ratings_anime_id ON ratings(anime_id);
CREATE INDEX idx_badges_user_id ON badges(user_id);
CREATE INDEX idx_watch_histories_user_id ON watch_histories(user_id);
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_conversation_participants_user_id ON conversation_participants(user_id);