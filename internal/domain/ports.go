package domain

import (
	"context"

	"github.com/google/uuid"
)

// ClientRepository is the port (interface) for persistence.
// The domain does not know about Postgres, HTTP or any infrastructure.
type ClientRepository interface {
	Create(ctx context.Context, client *Client) error
	GetByID(ctx context.Context, id uuid.UUID) (*Client, error)
	GetAll(ctx context.Context) ([]*Client, error)
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id uuid.UUID) error
}
