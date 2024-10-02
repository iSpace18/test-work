package models

import (
    "github.com/google/uuid"
)

type User struct {
    ID             uuid.UUID `db:"id"`
    Email          string    `db:"email"`
    RefreshTokenHash string   `db:"refresh_token_hash"`
    LastIP         string    `db:"last_ip"`
}