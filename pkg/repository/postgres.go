package repository

import (
	"database/sql"

	"vk/pkg/model"

	"github.com/jmoiron/sqlx"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (repo *PostgresRepository) GetDocument(url string) (*model.Document, error) {
	doc := &model.Document{}
	err := repo.db.Get(doc, "SELECT url, pub_date, fetch_time, text, first_fetch_time FROM documents WHERE url=$1", url)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}

	return doc, nil
}

func (repo *PostgresRepository) SaveDocument(doc *model.Document) error {
	_, err := repo.db.NamedExec(`INSERT INTO documents (url, pub_date, fetch_time, text, first_fetch_time) 
                                VALUES (:url, :pub_date, :fetch_time, :text, :first_fetch_time) 
                                ON CONFLICT (url) 
                                DO UPDATE SET pub_date = EXCLUDED.pub_date, 
                                              fetch_time = EXCLUDED.fetch_time,
                                              text = EXCLUDED.text, 
                                              first_fetch_time = EXCLUDED.first_fetch_time`, doc)
	return err
}

func (repo *PostgresRepository) LockDocument(url string) error {
	// Using PostgreSQL's advisory locks
	_, err := repo.db.Exec("SELECT pg_advisory_lock(hashtext($1))", url)
	return err
}

func (repo *PostgresRepository) UnlockDocument(url string) error {
	// Using PostgreSQL's advisory locks
	_, err := repo.db.Exec("SELECT pg_advisory_unlock(hashtext($1))", url)
	return err
}
