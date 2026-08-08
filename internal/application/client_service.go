package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/grok/crudclienteshex/internal/domain"
)

var (
	ErrClientNotFound   = errors.New("cliente no encontrado")
	ErrInvalidInput     = errors.New("datos de entrada inválidos")
	ErrEmailRequired    = errors.New("el email es obligatorio")
	ErrNombreRequired   = errors.New("el nombre es obligatorio")
)

// ClientService contains the use cases (application layer).
type ClientService struct {
	repo domain.ClientRepository
}

func NewClientService(repo domain.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) Create(ctx context.Context, nombre, email, telefono, direccion string) (*domain.Client, error) {
	nombre = strings.TrimSpace(nombre)
	email = strings.TrimSpace(email)
	telefono = strings.TrimSpace(telefono)
	direccion = strings.TrimSpace(direccion)

	if nombre == "" {
		return nil, ErrNombreRequired
	}
	if email == "" {
		return nil, ErrEmailRequired
	}

	client := domain.NewClient(nombre, email, telefono, direccion)
	if err := s.repo.Create(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *ClientService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, ErrClientNotFound
	}
	return client, nil
}

func (s *ClientService) GetAll(ctx context.Context) ([]*domain.Client, error) {
	return s.repo.GetAll(ctx)
}

func (s *ClientService) Update(ctx context.Context, id uuid.UUID, nombre, email, telefono, direccion string) (*domain.Client, error) {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, ErrClientNotFound
	}

	nombre = strings.TrimSpace(nombre)
	email = strings.TrimSpace(email)
	telefono = strings.TrimSpace(telefono)
	direccion = strings.TrimSpace(direccion)

	if nombre == "" {
		return nil, ErrNombreRequired
	}
	if email == "" {
		return nil, ErrEmailRequired
	}

	client.Nombre = nombre
	client.Email = email
	client.Telefono = telefono
	client.Direccion = direccion
	client.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *ClientService) Delete(ctx context.Context, id uuid.UUID) error {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if client == nil {
		return ErrClientNotFound
	}
	return s.repo.Delete(ctx, id)
}
