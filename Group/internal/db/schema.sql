-- ============================================
-- FORUM DATABASE SCHEMA
-- ============================================
-- Server-side rendered forum application
-- SQLite3 with foreign key enforcement
-- 
-- Phase 3 — Database Layer
-- Author: talo (Phase 3A)
-- ============================================

-- Enable foreign key support (disabled by default in SQLite)
PRAGMA foreign_keys = ON;

-- ============================================
-- USERS TABLE
-- ============================================
-- Stores registered forum accounts.
-- Email and username are unique to prevent duplicates.
-- Passwords are stored as bcrypt hashes, never plaintext.

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- SESSIONS TABLE
-- ============================================
-- Manages server-side session tokens stored in browser cookies.
-- Tokens are UUID strings with expiration timestamps.
-- CASCADE delete removes sessions when a user is deleted.

CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================
-- POSTS TABLE
-- ============================================
-- Forum threads created by registered users.
-- CASCADE delete removes posts when author is deleted.

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- ============================================
-- CATEGORIES TABLE
-- ============================================
-- Predefined forum sections (e.g., Technology, Sports, Gaming).
-- Seeded at application startup by repository layer.
-- Unique names prevent duplicate categories.

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- ============================================
-- POST_CATEGORIES TABLE
-- ============================================
-- Junction table for many-to-many relationship between
-- posts and categories. Composite primary key prevents
-- duplicate assignments.

CREATE TABLE IF NOT EXISTS post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    
    PRIMARY KEY (post_id, category_id),
    
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- ============================================
-- COMMENTS TABLE
-- ============================================
-- User replies attached to specific posts.
-- Dual foreign keys with CASCADE delete on both sides.

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================
-- REACTIONS TABLE
-- ============================================
-- Stores likes and dislikes on posts and comments.
-- Uses polymorphic design with CHECK constraint ensuring
-- exactly one target (post OR comment) per reaction.
-- 
-- Design rationale:
-- - UNIQUE(user_id, post_id) allows multiple NULL post_ids
--   (SQLite treats NULLs as distinct in UNIQUE constraints)
-- - Same logic applies to UNIQUE(user_id, comment_id)
-- - CHECK constraint enforces exactly one target type

CREATE TABLE IF NOT EXISTS reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    type TEXT NOT NULL CHECK(type IN ('like', 'dislike')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    
    CHECK (
        (post_id IS NOT NULL AND comment_id IS NULL) OR
        (post_id IS NULL AND comment_id IS NOT NULL)
    ),
    
    UNIQUE(user_id, post_id),
    UNIQUE(user_id, comment_id)
);

-- ============================================
-- PERFORMANCE INDEXES
-- ============================================
-- Optimize frequent query patterns identified in Phase 3 review:
-- - Post listings by author (profile page)
-- - Comment retrieval per post (post detail page)
-- - Reaction counts (post and comment cards)
-- - Session validation (every authenticated request)
-- - Session cleanup (expired session purging)

CREATE INDEX IF NOT EXISTS idx_posts_user ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_reactions_post ON reactions(post_id);
CREATE INDEX IF NOT EXISTS idx_reactions_comment ON reactions(comment_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);
