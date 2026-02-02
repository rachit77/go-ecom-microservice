package account

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close() error
	PutAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	// skip for pagination
	ListAccounts(ctx context.Context, skip uint, take uint) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) PutAccount(ctx context.Context, a Account) error {
	query := "INSERT INTO accounts(id, name) VALUES($1, $2)"
	_, err := r.db.ExecContext(ctx, query, a.ID, a.Name)
	return err
}

func (r *postgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	query := "SELECT id, name FROM accounts WHERE id = $1"
	row, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	a := &Account{}
	if err := row.Scan(&a.ID, &a.Name); err != nil {
		return nil, err
	}
	return a, err
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint, take uint) ([]Account, error) {
	query := "SELECT id, name from accounts ORDER BY id DESC OFFSET $1 LIMIT $2"
	rows, err := r.db.QueryContext(ctx, query, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []Account{}

	for rows.Next() {
		a := &Account{}
		if err := rows.Scan(&a.ID, &a.Name); err == nil {
			accounts = append(accounts, *a)
		}
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return accounts, nil
}
