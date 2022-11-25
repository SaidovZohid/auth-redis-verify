package repo

import "time"

type User struct {
	ID        int64     
	FirstName string    
	Email     string    
	CreatedAt time.Time 
}

type UserStorageI interface {
	Create(user *User) (*User, error)
}
