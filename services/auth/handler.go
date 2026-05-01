package auth

import "gorm.io/gorm"

type Handler struct {
	db        *gorm.DB
	jwtSecret []byte
}

func NewHandler(db *gorm.DB, jwtSecret string) *Handler {
	return &Handler{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}
