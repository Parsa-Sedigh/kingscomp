package matchmaking

import (
	"context"
	_ "embed"
	"errors"
	"github.com/Parsa-Sedigh/kingscomp/internal/entity"
	"github.com/Parsa-Sedigh/kingscomp/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/rueidis"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	ErrBadRedisResponse = errors.New("bad redis response")
	ErrTimeout          = errors.New("lobby queue timeout")
)

type MatchMaking interface {
	Join(ctx context.Context, userID int64, timeout time.Duration) (entity.Lobby, error)
	Leave(ctx context.Context, userID int64) error
}

//go:embed matchmaking.lua
var matchMakingScript string

var _ MatchMaking = &RedisMatchMaking{}

type RedisMatchMaking struct {
	client            rueidis.Client
	matchMakingScript *rueidis.Lua
	lobby             repository.LobbyRepository
}

func NewRedisMatchMaking(client rueidis.Client, lobby repository.LobbyRepository) *RedisMatchMaking {
	script := rueidis.NewLuaScript(matchMakingScript)

	return &RedisMatchMaking{
		client:            client,
		matchMakingScript: script,
		lobby:             lobby,
	}
}

type JoinLobbyPubSubResponse struct {
	err     error
	lobbyID string
}

// Join blocks until the userID joins the lobby or the func time
func (r RedisMatchMaking) Join(ctx context.Context, userID int64, timeout time.Duration) (entity.Lobby, error) {
	waitingLobbyCtx, lobbyCtxCancel := context.WithTimeout(ctx, time.Minute*2)
	defer lobbyCtxCancel()

	responseChannel := make(chan JoinLobbyPubSubResponse, 1)

	go r.client.Receive(waitingLobbyCtx, r.client.B().Subscribe().Channel("matchmaking").Build(), func(msg rueidis.PubSubMessage) {
		message := strings.Split(msg.Message, ":")
		lobbyId := message[0]
		users := lo.Map(strings.Split(message[1], ","), func(item string, _ int) int64 {
			id, _ := strconv.ParseInt(item, 10, 64)

			return id
		})

		if !slices.Contains(users, userID) {
			return
		}

		responseChannel <- JoinLobbyPubSubResponse{
			lobbyID: lobbyId,
		}
	})

	resp, err := r.matchMakingScript.Exec(ctx, r.client,
		[]string{"matchmaking", "matchmaking", strconv.FormatInt(userID, 10)},
		[]string{"4",
			strconv.FormatInt(time.Now().Add(-time.Minute*2).Unix(), 10),
			uuid.New().String(),
			strconv.FormatInt(userID, 10),
			strconv.FormatInt(time.Now().Unix(), 10),
		},
	).ToArray()
	if err != nil {
		logrus.WithError(err).Errorln("couldn't join the match making")
	}

	// current user moved into a queue. We must listen to the pub/sub
	if len(resp) == 1 {
		select {
		case pubSubRes := <-responseChannel:
			return r.lobby.Get(ctx, entity.NewID("lobby", pubSubRes.lobbyID))

		case <-waitingLobbyCtx.Done():
			return entity.Lobby{}, err
		}
	}

	// a lobby just created
	if len(resp) == 3 {
		lobbyID, _ := resp[1].ToString()

		return r.lobby.Get(ctx, entity.NewID("lobby", lobbyID))
	}

	return entity.Lobby{}, ErrBadRedisResponse
}

func (r RedisMatchMaking) Leave(ctx context.Context, userID int64) error {
	return nil
}
