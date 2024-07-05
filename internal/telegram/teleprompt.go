package telegram

import (
	"gopkg.in/telebot.v3"
	"sync"
	"time"
)

type Prompt struct {
	TeleCtx telebot.Context
}

type TelePrompt struct {
	accountPrompts sync.Map
}

func NewTelePrompt() *TelePrompt {
	return &TelePrompt{}
}

func (t *TelePrompt) Register(userID int64) <-chan Prompt {
	c := make(chan Prompt, 1)

	if preChannel, loaded := t.accountPrompts.LoadAndDelete(userID); loaded {
		close(preChannel.(chan Prompt))
	}

	t.accountPrompts.Store(userID, c)

	return c
}

func (t *TelePrompt) AsMessage(userID int64, timeout time.Duration) (*telebot.Message, bool) {
	c := t.Register(userID)

	select {
	case val := <-c:
		return val.TeleCtx.Message(), false

	case <-time.After(timeout):
		return nil, true
	}
}

func (t *TelePrompt) Dispatch(userID int64, c telebot.Context) bool {
	ch, loaded := t.accountPrompts.LoadAndDelete(userID)
	if !loaded {
		return false
	}

	select {
	case ch.(chan Prompt) <- Prompt{TeleCtx: c}:

	default:
		return false
	}

	return true
}
