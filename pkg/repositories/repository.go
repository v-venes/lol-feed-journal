package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Repository[T any] struct {
	db *sqlx.DB
}

func NewRepository[T any](db *sqlx.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) GetAll(query string, args ...interface{}) ([]T, error) {
	var entities []T
	if err := r.db.Select(&entities, query, args...); err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *Repository[T]) SaveMany(query string, data []T) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = r.db.NamedExec(query, data)

	if err != nil {
		return err
	}

	return tx.Commit()
}
