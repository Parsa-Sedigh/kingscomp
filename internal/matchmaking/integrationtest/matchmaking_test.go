package integrationtest

import (
	"context"
	"fmt"
	"github.com/Parsa-Sedigh/kingscomp/internal/matchmaking"
	"github.com/Parsa-Sedigh/kingscomp/internal/repository"
	"github.com/Parsa-Sedigh/kingscomp/internal/repository/redis"
	"github.com/redis/rueidis"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestMatchmaking_Join(t *testing.T) {
	ctx := context.Background()
	timeout := time.Second * 10
	redisClient, err := redis.NewRedisClient(fmt.Sprintf("localhost:%s", redisPort))
	assert.NoError(t, err)

	mm := matchmaking.NewRedisMatchMaking(redisClient, repository.NewLobbyRedisRepository(redisClient))

	var wg sync.WaitGroup

	testJoin := func(id int64) {
		wg.Add(1)

		go func() {
			lobby, err := mm.Join(ctx, id, timeout)
			assert.NoError(t, err)
			assert.NotEqual(t, "", lobby.ID)

			wg.Done()
		}()
	}

	testJoin(10)
	testJoin(11)
	testJoin(12)
	testJoin(13)

	<-time.After(time.Millisecond * 500)

	assert.Equal(t, 4, zCount(t, redisClient, "matchmaking"))

	lobby, err := mm.Join(ctx, 14, timeout)
	assert.NoError(t, err)
	assert.NotEqual(t, "", lobby.ID)
	wg.Wait()

	// check if the lobby has been created
	allKeys := keys(t, redisClient, "*")
	assert.Len(t, allKeys, 1)

	lobbyKey := allKeys[0]
	assert.Contains(t, lobbyKey, "lobby:")

	redisClient.Do(ctx, redisClient.B().JsonGet().Key(lobbyKey).Path(".").Build()).ToString()
}

type LobbyCounter struct {
	sync.Mutex
	counter map[string]int
}

func (l *LobbyCounter) Incr(lobbyID string) {
	l.Lock()
	defer l.Unlock()

	l.counter[lobbyID]++
}

func newLobbyCounter() *LobbyCounter {
	return &LobbyCounter{
		Mutex:   sync.Mutex{},
		counter: make(map[string]int),
	}
}

func TestMatchmaking_JoinWithManyLobbies(t *testing.T) {
	ctx := context.Background()
	timeout := time.Second * 10
	redisClient, err := redis.NewRedisClient(fmt.Sprintf("localhost:%s", redisPort))
	assert.NoError(t, err)

	mm := matchmaking.NewRedisMatchMaking(redisClient, repository.NewLobbyRedisRepository(redisClient))
	lobbyCounter := newLobbyCounter()

	var wg sync.WaitGroup

	testJoin := func(id int64) {
		wg.Add(1)

		go func() {
			lobby, err := mm.Join(ctx, id, timeout)
			assert.NoError(t, err)
			assert.NotEqual(t, "", lobby.ID)

			lobbyCounter.Incr(lobby.ID)

			wg.Done()
		}()
	}

	for i := 0; i < 100; i++ {
		testJoin(int64(i) + 1)
	}

	wg.Wait()

	assert.Len(t, lobbyCounter.counter, 20)

	for lobbyID, count := range lobbyCounter.counter {

	}
}

func zCount(t *testing.T, redisClient rueidis.Client, key string) int64 {
	count, err := redisClient.Do(context.Background(),
		redisClient.B().Zcount().Key(key).Min("-inf").Max("+inf").Build()).
		ToInt64()
	assert.NoError(t, err)

	return count
}

func keys(t *testing.T, redisClient rueidis.Client, pattern string) []string {
	items, err := redisClient.Do(context.Background(),
		redisClient.B().Keys().Pattern(pattern).Build()).
		AsStrSlice()
	assert.NoError(t, err)

	return items
}
