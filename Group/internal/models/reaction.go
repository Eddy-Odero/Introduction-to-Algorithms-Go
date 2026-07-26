package models

import "database/sql"

type ReactionRepository struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

func (r *ReactionRepository) SetPostReaction(userID, postID int, reactionType string) error {
	var existing string
	err := r.db.QueryRow(
		`SELECT type FROM reactions WHERE user_id = ? AND post_id = ?`,
		userID, postID,
	).Scan(&existing)

	switch {
	case err == sql.ErrNoRows:
		_, err = r.db.Exec(
			`INSERT INTO reactions (user_id, post_id, type) VALUES (?, ?, ?)`,
			userID, postID, reactionType,
		)
		return err
	case err != nil:
		return err
	case existing == reactionType:
		_, err = r.db.Exec(
			`DELETE FROM reactions WHERE user_id = ? AND post_id = ?`,
			userID, postID,
		)
		return err
	default:
		_, err = r.db.Exec(
			`UPDATE reactions SET type = ? WHERE user_id = ? AND post_id = ?`,
			reactionType, userID, postID,
		)
		return err
	}
}

func (r *ReactionRepository) SetCommentReaction(userID, commentID int, reactionType string) error {
	var existing string
	err := r.db.QueryRow(
		`SELECT type FROM reactions WHERE user_id = ? AND comment_id = ?`,
		userID, commentID,
	).Scan(&existing)

	switch {
	case err == sql.ErrNoRows:
		_, err = r.db.Exec(
			`INSERT INTO reactions (user_id, comment_id, type) VALUES (?, ?, ?)`,
			userID, commentID, reactionType,
		)
		return err
	case err != nil:
		return err
	case existing == reactionType:
		_, err = r.db.Exec(
			`DELETE FROM reactions WHERE user_id = ? AND comment_id = ?`,
			userID, commentID,
		)
		return err
	default:
		_, err = r.db.Exec(
			`UPDATE reactions SET type = ? WHERE user_id = ? AND comment_id = ?`,
			reactionType, userID, commentID,
		)
		return err
	}
}

func (r *ReactionRepository) GetUserPostVotes(userID int) (map[int]string, error) {
	votes := make(map[int]string)
	rows, err := r.db.Query(
		`SELECT post_id, type FROM reactions WHERE user_id = ? AND post_id IS NOT NULL`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var t string
		if err := rows.Scan(&id, &t); err != nil {
			return nil, err
		}
		votes[id] = t
	}
	return votes, rows.Err()
}

// CountLikesGivenByUserID returns how many "like" reactions this user has
// made across posts and comments combined — used on the profile page.
func (r *ReactionRepository) CountLikesGivenByUserID(userID int) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM reactions WHERE user_id = ? AND type = 'like'`,
		userID,
	).Scan(&count)
	return count, err
}