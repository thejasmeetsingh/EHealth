package models

import (
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/internal/database"
)

type user struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsEndUser bool      `json:"is_end_user"`
}

func DatabaseUserToUser(dbUser database.User) user {
	return user{
		ID:        dbUser.ID,
		Name:      dbUser.Name.String,
		Email:     dbUser.Email,
		IsEndUser: dbUser.IsEndUser,
	}
}
