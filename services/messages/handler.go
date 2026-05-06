package messages

import "gorm.io/gorm"

type Handler struct {
	db        *gorm.DB
	hub       *Hub
	jwtSecret []byte
}

func NewHandler(db *gorm.DB, hub *Hub, jwtSecret []byte) *Handler {
	return &Handler{
		db:        db,
		hub:       hub,
		jwtSecret: jwtSecret,
	}
}
