package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/grok/crudclienteshex/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClientRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresClientRepository(pool *pgxpool.Pool) *PostgresClientRepository {
	return &PostgresClientRepository{pool: pool}
}

func (r *PostgresClientRepository) Create(ctx context.Context, client *domain.Client) error {
	query := `
		INSERT INTO clients (id, nombre, email, telefono, direccion, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(ctx, query,
		client.ID,
		client.Nombre,
		client.Email,
		client.Telefono,
		client.Direccion,
		client.CreatedAt,
		client.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	return nil
}

func (r *PostgresClientRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	query := `
		SELECT id, nombre, email, telefono, direccion, created_at, updated_at
		FROM clients
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var c domain.Client
	err := row.Scan(
		&c.ID,
		&c.Nombre,
		&c.Email,
		&c.Telefono,
		&c.Direccion,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get client by id: %w", err)
	}
	return &c, nil
}

func (r *PostgresClientRepository) GetAll(ctx context.Context) ([]*domain.Client, error) {
	query := `
		SELECT id, nombre, email, telefono, direccion, created_at, updated_at
		FROM clients
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all clients: %w", err)
	}
	defer rows.Close()

	var clients []*domain.Client
	for rows.Next() {
		var c domain.Client
		if err := rows.Scan(
			&c.ID,
			&c.Nombre,
			&c.Email,
			&c.Telefono,
			&c.Direccion,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan client: %w", err)
		}
		clients = append(clients, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return clients, nil
}

func (r *PostgresClientRepository) Update(ctx context.Context, client *domain.Client) error {
	query := `
		UPDATE clients
		SET nombre = $2, email = $3, telefono = $4, direccion = $5, updated_at = $6
		WHERE id = $1
	`
	cmd, err := r.pool.Exec(ctx, query,
		client.ID,
		client.Nombre,
		client.Email,
		client.Telefono,
		client.Direccion,
		client.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update client: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return nil // no rows updated
	}
	return nil
}

func (r *PostgresClientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM clients WHERE id = $1`
	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete client: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return nil // or return not found
	}
	return nil
}

// Ensure table exists (simple migration helper)
func (r *PostgresClientRepository) EnsureSchema(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS clients (
			id UUID PRIMARY KEY,
			nombre VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			telefono VARCHAR(50),
			direccion TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_clients_email ON clients(email);
	`
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("ensure schema: %w", err)
	}
	return nil
}
