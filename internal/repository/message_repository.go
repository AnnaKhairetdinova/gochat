package repository

import (
	"context"

	"github.com/AnnaKhairetdinova/gochat/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository interface {
	Save(ctx context.Context, msg domain.Message) error
	GetHistory(ctx context.Context, room string, limit int) ([]domain.Message, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) MessageRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Save(ctx context.Context, msg domain.Message) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	queryMessage := `INSERT INTO messages (uuid, room, username, text, type, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(ctx, queryMessage, msg.UUID, msg.Room, msg.Username, msg.Text, msg.Type, msg.CreatedAt)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *postgresRepository) GetHistory(ctx context.Context, room string, limit int) ([]domain.Message, error) {
	queryMessage := `SELECT * FROM (SELECT uuid, room, username, text, type, created_at FROM messages WHERE room = $1 ORDER BY created_at DESC LIMIT $2) AS recent ORDER BY created_at ASC`

	rows, err := r.db.Query(ctx, queryMessage, room, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []domain.Message

	for rows.Next() {
		var message domain.Message
		err := rows.Scan(&message.UUID, &message.Room, &message.Username, &message.Text, &message.Type, &message.CreatedAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if messages == nil {
		messages = []domain.Message{}
	}

	return messages, nil
}
