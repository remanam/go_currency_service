package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/remanam/go_currency_service/auth_service/internal/domain"
)

type PostgresUserRepo struct {
	conn *pgx.Conn
}

func NewPostgresUserRepo(conn *pgx.Conn) *PostgresUserRepo {
	return &PostgresUserRepo{conn: conn}
}
func (r *PostgresUserRepo) Create(user *domain.User) (int32, error) {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`

	var id int32
	err := r.conn.QueryRow(context.Background(), query, user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresUserRepo) GetByUsername(username string) (*domain.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE username = $1`
	var user domain.User
	err := r.conn.QueryRow(context.Background(), query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepo) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = $1`

	var user domain.User
	err := r.conn.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
