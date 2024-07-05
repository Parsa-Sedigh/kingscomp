package repository

import (
	"context"
	"errors"
	"github.com/Parsa-Sedigh/kingscomp/internal/entity"
)

// we want to have the same error when dealing with DB or redis or ... .

var (
	ErrNotFound = errors.New("entity not found")
)

type CommonBehaviour[T entity.Entity] interface {
	Get(ctx context.Context, id entity.ID) (T, error)
	Save(ctx context.Context, ent entity.Entity) error
	//UpdateField(ctx context.Context, id entity.ID, field string, val any) error
}

//go:generate mockery --name AccountRepository
type AccountRepository interface {
	CommonBehaviour[entity.Account]
}

//go:generate mockery --name LobbyRepository
type LobbyRepository interface {
	CommonBehaviour[entity.Lobby]
}
