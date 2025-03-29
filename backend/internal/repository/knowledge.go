package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"support/internal/models"
)

type KnowledgeRepo struct {
	db *sqlx.DB
}

func NewKnowledgeRepo(db *sqlx.DB) *KnowledgeRepo {
	return &KnowledgeRepo{db: db}
}

func (r *KnowledgeRepo) Get(ctx context.Context, intent string) (models.Knowledge, error) {
	query, args, err := sq.Select("*").
		From("knowledge").
		Where(sq.Eq{"intent": intent}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return models.Knowledge{}, err
	}

	var knowledge models.Knowledge
	err = r.db.GetContext(ctx, &knowledge, query, args...)
	if err != nil {
		return models.Knowledge{}, err
	}

	return knowledge, nil
}

func (r *KnowledgeRepo) Exist(ctx context.Context, intent string) bool {
	query, args, err := sq.Select("COUNT(*) > 0").
		From("knowledge").
		Where(sq.Eq{"intent": intent}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false
	}

	var exists bool
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
		return false
	}
	return exists
}
