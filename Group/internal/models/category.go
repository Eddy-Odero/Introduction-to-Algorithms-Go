package models

import (
	"database/sql"
	"fmt"
	"strings"  
)

// Category represents a forum topic tag
type Category struct {
	ID   int
	Name string
}

// CategoryRepository handles all category database operations
type CategoryRepository struct {
	DB *sql.DB
}

// NewCategoryRepository creates a repository instance
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

// GetAll returns every category in alphabetical order
func (r *CategoryRepository) GetAll() ([]Category, error) {
	rows, err := r.DB.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("GetAll categories query failed: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("GetAll scan failed: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

// GetByName looks up a single category by its name
func (r *CategoryRepository) GetByName(name string) (*Category, error) {
	var c Category
	err := r.DB.QueryRow(
		"SELECT id, name FROM categories WHERE name = ?",
		name,
	).Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetByName failed: %w", err)
	}
	return &c, nil
}

// GetByPostID returns all categories linked to a specific post
func (r *CategoryRepository) GetByPostID(postID int) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT c.name 
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
		ORDER BY c.name
	`, postID)
	if err != nil {
		return nil, fmt.Errorf("GetByPostID query failed: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("GetByPostID scan failed: %w", err)
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// LinkToPost creates the many-to-many relationship between a post and categories
func (r *CategoryRepository) LinkToPost(postID int, categoryNames []string) error {
	for _, name := range categoryNames {
		name = strings.TrimSpace(name)  // ADDED: prevent "Tech " vs "Tech" duplicates
		if name == "" {
			continue
		}
		
		// Find category ID
		var catID int
		err := r.DB.QueryRow(
			"SELECT id FROM categories WHERE name = ?",
			name,
		).Scan(&catID)
		if err == sql.ErrNoRows {
			// Auto-create unknown categories
			res, err := r.DB.Exec(
				"INSERT INTO categories (name) VALUES (?)",
				name,
			)
			if err != nil {
				return fmt.Errorf("auto-create category failed: %w", err)
			}
			lastID, _ := res.LastInsertId()
			catID = int(lastID)
		} else if err != nil {
			return fmt.Errorf("lookup category failed: %w", err)
		}

		// Link post to category
		_, err = r.DB.Exec(
			"INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)",
			postID, catID,
		)
		if err != nil {
			return fmt.Errorf("link category failed: %w", err)
		}
	}
	return nil
}