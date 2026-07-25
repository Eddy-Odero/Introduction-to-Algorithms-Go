package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	ID           int
	UserID       int
	Title        string
	Author       string
	Excerpt      string
	Content      string
	ImageURL     string
	Categories   []string
	Likes        int
	Dislikes     int
	CommentCount int
	CreatedAt    time.Time
}

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) CreatePost(userID int, title, content string) (int64, error) {
	result, err := r.DB.Exec(
		"INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		userID, title, content,
	)
	if err != nil {
		return 0, fmt.Errorf("create post failed: %w", err)
	}
	return result.LastInsertId()
}

func (r *PostRepository) GetAll() ([]Post, error) {
	rows, err := r.DB.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("get all posts failed: %w", err)
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *PostRepository) GetByID(id int) (*Post, error) {
	post := &Post{}
	err := r.DB.QueryRow(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Author)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get post by id failed: %w", err)
	}

	post.Excerpt = makeExcerpt(post.Content)
	post.Categories, _ = r.getCategories(post.ID)
	post.Likes, post.Dislikes, _ = r.getReactionCounts(post.ID)
	post.CommentCount, _ = r.getCommentCount(post.ID)

	return post, nil
}

func (r *PostRepository) GetByCategory(categoryName string) ([]Post, error) {
	rows, err := r.DB.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		JOIN categories c ON pc.category_id = c.id
		WHERE c.name = ?
		ORDER BY p.created_at DESC
	`, categoryName)
	if err != nil {
		return nil, fmt.Errorf("get posts by category failed: %w", err)
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *PostRepository) GetByUserID(userID int) ([]Post, error) {
	rows, err := r.DB.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("get posts by user failed: %w", err)
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *PostRepository) GetLikedByUserID(userID int) ([]Post, error) {
	rows, err := r.DB.Query(`
		SELECT DISTINCT p.id, p.user_id, p.title, p.content, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN reactions r ON p.id = r.post_id
		WHERE r.user_id = ? AND r.type = 'like'
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("get liked posts failed: %w", err)
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *PostRepository) scanPosts(rows *sql.Rows) ([]Post, error) {
	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.Author); err != nil {
			return nil, fmt.Errorf("scan post failed: %w", err)
		}
		p.Excerpt = makeExcerpt(p.Content)
		p.Categories, _ = r.getCategories(p.ID)
		p.Likes, p.Dislikes, _ = r.getReactionCounts(p.ID)
		p.CommentCount, _ = r.getCommentCount(p.ID)
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *PostRepository) getCategories(postID int) ([]string, error) {
	catRepo := NewCategoryRepository(r.DB)
	return catRepo.GetByPostID(postID)
}

func (r *PostRepository) getReactionCounts(postID int) (likes, dislikes int, err error) {
	err = r.DB.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'like' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'dislike' THEN 1 ELSE 0 END), 0)
		FROM reactions
		WHERE post_id = ?
	`, postID).Scan(&likes, &dislikes)
	if err != nil {
		return 0, 0, nil
	}
	return likes, dislikes, nil
}

func (r *PostRepository) getCommentCount(postID int) (int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", postID).Scan(&count)
	if err != nil {
		return 0, nil
	}
	return count, nil
}

func makeExcerpt(content string) string {
	if len(content) <= 150 {
		return content
	}
	return content[:150] + "..."
}