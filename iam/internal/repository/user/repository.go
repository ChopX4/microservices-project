package user

import "github.com/jackc/pgx/v5/pgxpool"

type repository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}
