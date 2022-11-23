package postgres

import (
	"github.com/SaidovZohid/auth-redis-verify/storage/repo"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) Create(user *repo.User) (*repo.User, error) {
	tr, err := ur.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tr.Rollback()
	query := `
		INSERT INTO users (
			first_name,
			email
		) VALUES ($1, $2)
		RETURNING id, created_at	
	`
	var u repo.User
	err = tr.QueryRow(
		query,
		user.FirstName,
		user.Email,
	).Scan(
		&u.ID,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	tr.Commit()
	return &u, nil
}
