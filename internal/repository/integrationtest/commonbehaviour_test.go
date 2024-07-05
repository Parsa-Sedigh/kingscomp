package integrationtest

import (
	"context"
	"fmt"
	"github.com/Parsa-Sedigh/kingscomp/internal/entity"
	"github.com/Parsa-Sedigh/kingscomp/internal/repository"
	"github.com/Parsa-Sedigh/kingscomp/internal/repository/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testType struct {
	ID   string
	Name string
}

func (t testType) EntityID() entity.ID {
	return entity.NewID("testType", t.ID)
}

func TestCommonBehaviourSetAndGet(t *testing.T) {
	redisClient, err := redis.NewRedisClient(fmt.Sprintf("localhost:%s", redisPort))
	assert.NoError(t, err)

	fmt.Println("Redis is now connected")

	cb := repository.NewRedisCommonBehavior[testType](redisClient)
	ctx := context.Background()

	err = cb.Save(ctx, &testType{
		ID:   "12",
		Name: "Parsa Sedigh",
	})
	assert.NoError(t, err)

	err = cb.Save(ctx, &testType{
		ID:   "13",
		Name: "Mehdi",
	})
	assert.NoError(t, err)

	val, err := cb.Get(ctx, entity.NewID("testType", "12"))
	assert.NoError(t, err)
	assert.Equal(t, "Parsa Sedigh", val.Name)
	assert.Equal(t, "12", val.ID)

	val, err = cb.Get(ctx, entity.NewID("testType", "13"))
	assert.NoError(t, err)
	assert.Equal(t, "Mehdi", val.Name)
	assert.Equal(t, "13", val.ID)

	err = cb.Save(ctx, &testType{
		ID:   "13",
		Name: "Yasin",
	})
	assert.NoError(t, err)

	val, err = cb.Get(ctx, entity.NewID("testType", "13"))
	assert.NoError(t, err)
	assert.Equal(t, "yasin", val.Name)
	assert.Equal(t, "13", val.ID)

	val, err = cb.Get(ctx, entity.NewID("testType", "14"))
	assert.ErrorIs(t, repository.ErrNotFound, err)

	redisClient.Close()
}
