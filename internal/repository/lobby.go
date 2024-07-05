package repository

import (
	"github.com/Parsa-Sedigh/kingscomp/internal/entity"
	"github.com/redis/rueidis"
)

var _ LobbyRepository = &LobbyRedisRepository{}

type LobbyRedisRepository struct {
	*RedisCommonBehavior[entity.Lobby]
}

func NewLobbyRedisRepository(client rueidis.Client) *LobbyRedisRepository {
	return &LobbyRedisRepository{
		RedisCommonBehavior: NewRedisCommonBehavior[entity.Lobby](client),
	}
}
