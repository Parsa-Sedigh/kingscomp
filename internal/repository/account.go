package repository

import (
	"github.com/Parsa-Sedigh/kingscomp/internal/entity"
	"github.com/redis/rueidis"
)

// make sure AccountRedisRepository implements AccountRepository interface
var _ AccountRepository = &AccountRedisRepository{}

type AccountRedisRepository struct {
	*RedisCommonBehavior[entity.Account]
}

func NewAccountRedisRepository(client rueidis.Client) *AccountRedisRepository {
	return &AccountRedisRepository{
		RedisCommonBehavior: NewRedisCommonBehavior[entity.Account](client),
	}
}
