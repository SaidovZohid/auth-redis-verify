package storage

import (
	"github.com/SaidovZohid/auth-redis-verify/storage/postgres"
	"github.com/SaidovZohid/auth-redis-verify/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
}

type StoragePg struct {
	userRepo repo.UserStorageI
}

func NewStorage(db *sqlx.DB) StorageI {
	return &StoragePg{
		userRepo: postgres.NewUser(db),
	}
}

func (s *StoragePg) User() repo.UserStorageI {
	return s.userRepo
}