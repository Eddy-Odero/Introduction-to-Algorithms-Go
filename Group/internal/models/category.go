package models

import "database/sql"

type Category struct {
	ID   int64
	Name string
}

type CategoryStore struct {
	DB *sql.DB
}

// All returns every category, ordered by name — used to render the
// category picker on new-post and the filter row on the feed.
func (s *CategoryStore) All() ([]Category, error) {
	rows, err := s.DB.Query(`SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

// IDsByNames resolves category names (as submitted by the new-post form)
// to their IDs, silently skipping any name that doesn't match a real
// category rather than erroring — the checkbox values are trusted to
// match seeded names, but this stays defensive against a stale/edited form.
func (s *CategoryStore) IDsByNames(names []string) ([]int64, error) {
	var ids []int64
	for _, name := range names {
		var id int64
		err := s.DB.QueryRow(`SELECT id FROM categories WHERE name = ?`, name).Scan(&id)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
