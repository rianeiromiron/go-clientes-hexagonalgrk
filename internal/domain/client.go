package domain

import (
	"time"

	"github.com/google/uuid"
)

// Client is the core domain entity.
type Client struct {
	ID        uuid.UUID `json:"id"`
	Nombre    string    `json:"nombre"`
	Email     string    `json:"email"`
	Telefono  string    `json:"telefono"`
	Direccion string    `json:"direccion"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewClient creates a new Client with generated ID and timestamps.
func NewClient(nombre, email, telefono, direccion string) *Client {
	now := time.Now().UTC()
	return &Client{
		ID:        uuid.New(),
		Nombre:    nombre,
		Email:     email,
		Telefono:  telefono,
		Direccion: direccion,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
