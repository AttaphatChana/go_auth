package persistence

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "test_auth/application/ports"
	"test_auth/domain"
)

type PostgresUserRepo struct{ Pool *pgxpool.Pool }

func (r *PostgresUserRepo) Save(u *domain.User) error {
	_, err := r.Pool.Exec(context.Background(),
		`INSERT INTO users(id,username,password) VALUES($1,$2,$3)`,
		u.ID, u.Username, u.Password)
	return err
}
func (r *PostgresUserRepo) FindByUsername(username string) (*domain.User, error) {
	row := r.Pool.QueryRow(context.Background(),
		`SELECT id,username,password FROM users WHERE username=$1`, username)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Username, &u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}
