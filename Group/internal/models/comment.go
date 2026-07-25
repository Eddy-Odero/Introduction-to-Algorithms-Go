package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Comment represents a user reply to a forum post.
type Comment struct {
	ID        int64
	PostID    int64
	UserID    int64
	Username  string // Joined from users table
	Content   string
	CreatedAt time.Time
}

// CommentRepository provides database operations for comments.
type CommentRepository struct {
	DB *sql.DB
}

// NewCommentRepository creates a new CommentRepository.
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// CreateComment inserts a new comment. Returns the created comment with metadata.
func (cr *CommentRepository) CreateComment(postID, userID int64, content string) (*Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}

	result, err := cr.DB.Exec(
		`INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`,
		postID, userID, content,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return cr.GetCommentByID(id)
}

// GetCommentByID retrieves a comment by ID with author username.
func (cr *CommentRepository) GetCommentByID(id int64) (*Comment, error) {
	comment := &Comment{}
	err := cr.DB.QueryRow(
		`SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at 
		 FROM comments c
		 JOIN users u ON c.user_id = u.id
		 WHERE c.id = ?`,
		id,
	).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Username, &comment.Content, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// GetCommentsByPost retrieves all comments for a post, oldest first.
func (cr *CommentRepository) GetCommentsByPost(postID int64) ([]Comment, error) {
	rows, err := cr.DB.Query(
		`SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at 
		 FROM comments c
		 JOIN users u ON c.user_id = u.id
		 WHERE c.post_id = ?
		 ORDER BY c.created_at ASC`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

// CountCommentsByPost returns the number of comments on a post.
func (cr *CommentRepository) CountCommentsByPost(postID int64) (int64, error) {
	var count int64
	err := cr.DB.QueryRow(
		`SELECT COUNT(*) FROM comments WHERE post_id = ?`,
		postID,
	).Scan(&count)
	return count, err
}
// CountCommentsByUserID returns total comments written by a user
func (cr *CommentRepository) CountCommentsByUserID(userID int64) (int64, error) {
	var count int64
	err := cr.DB.QueryRow(
		`SELECT COUNT(*) FROM comments WHERE user_id = ?`,
		userID,
	).Scan(&count)
	return count, err
}
// DeleteComment removes a comment by ID.
func (cr *CommentRepository) DeleteComment(id int64) error {
	_, err := cr.DB.Exec(`DELETE FROM comments WHERE id = ?`, id)
	return err
}
